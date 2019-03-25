package lndclient

import (
	"context"
	"time"

	"github.com/lightningnetwork/lnd/channeldb"
	"google.golang.org/grpc/codes"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/lntypes"
)

// RouterClient exposes payment functionality.
type RouterClient interface {
	SendPayment(ctx context.Context, invoice string,
		maxFee btcutil.Amount, maxCltv *int32, outgoingChannel *uint64,
		timeout time.Duration) error

	LookupPayment(ctx context.Context, hash lntypes.Hash) (
		chan PaymentStatus, chan error, error)
}

// PaymentStatus describe the state of a payment.
type PaymentStatus struct {
	State    routerrpc.PaymentState
	Preimage lntypes.Preimage
	Fee      lnwire.MilliSatoshi
	Route    []uint64
}

type routerClient struct {
	client       routerrpc.RouterClient
	routerKitMac serializedMacaroon
}

func newRouterClient(conn *grpc.ClientConn,
	routerKitMac serializedMacaroon) *routerClient {

	return &routerClient{
		client:       routerrpc.NewRouterClient(conn),
		routerKitMac: routerKitMac,
	}
}

func (r *routerClient) SendPayment(ctx context.Context, invoice string,
	maxFee btcutil.Amount, maxCltv *int32, outgoingChannel *uint64,
	timeout time.Duration) error {

	rpcCtx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()

	rpcCtx = r.routerKitMac.WithMacaroonAuth(rpcCtx)
	req := &routerrpc.SendPaymentRequest{
		FeeLimit:       int64(maxFee),
		PaymentRequest: invoice,
		TimeoutSeconds: int32(timeout.Seconds()),
	}
	if maxCltv != nil {
		req.CltvLimit = *maxCltv
	}
	if outgoingChannel != nil {
		req.OutgoingChanId = *outgoingChannel
	}

	_, err := r.client.SendPayment(rpcCtx, req)

	return err
}

func (r *routerClient) LookupPayment(ctx context.Context,
	hash lntypes.Hash) (chan PaymentStatus, chan error, error) {

	ctx = r.routerKitMac.WithMacaroonAuth(ctx)
	stream, err := r.client.LookupPayment(
		ctx, &routerrpc.LookupPaymentRequest{
			PaymentHash: hash[:],
		},
	)
	if err != nil {
		return nil, nil, err
	}

	statusChan := make(chan PaymentStatus)
	errorChan := make(chan error, 1)
	go func() {
		for {
			rpcStatus, err := stream.Recv()
			if err != nil {
				if status.Convert(err).Code() ==
					codes.NotFound {

					err = channeldb.ErrPaymentNotInitiated
				}

				errorChan <- err
				return
			}

			status, err := unmarshallPaymentStatus(rpcStatus)
			if err != nil {
				errorChan <- err
				return
			}

			select {
			case statusChan <- *status:
			case <-ctx.Done():
				return
			}
		}
	}()

	return statusChan, errorChan, nil
}

func unmarshallPaymentStatus(rpcStatus *routerrpc.PaymentStatus) (
	*PaymentStatus, error) {

	status := PaymentStatus{
		State: rpcStatus.State,
	}

	if status.State == routerrpc.PaymentState_SUCCEEDED {
		preimage, err := lntypes.MakePreimage(
			rpcStatus.Preimage,
		)
		if err != nil {
			return nil, err
		}
		status.Preimage = preimage

		status.Fee = lnwire.MilliSatoshi(
			rpcStatus.Route.TotalFeesMsat,
		)

		status.Route = make([]uint64, len(rpcStatus.Route.Hops))
		for i, hop := range rpcStatus.Route.Hops {
			status.Route[i] = hop.ChanId
		}
	}

	return &status, nil
}
