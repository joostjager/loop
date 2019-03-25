package test

import (
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/lightninglabs/loop/lndclient"
	"github.com/lightningnetwork/lnd/lntypes"
	"golang.org/x/net/context"
)

type mockRouter struct {
	lnd *LndMockServices
}

func (r *mockRouter) SendPayment(ctx context.Context, invoice string,
	maxFee btcutil.Amount, maxCltv *int32, outgoingChannel *uint64,
	timeout time.Duration) error {

	r.lnd.RouterSendPaymentChannel <- RouterPaymentChannelMessage{
		Invoice:         invoice,
		MaxFee:          maxFee,
		MaxCltv:         maxCltv,
		OutgoingChannel: outgoingChannel,
	}

	return nil
}

func (r *mockRouter) LookupPayment(ctx context.Context,
	hash lntypes.Hash) (chan lndclient.PaymentStatus, chan error, error) {

	statusChan := make(chan lndclient.PaymentStatus)
	errorChan := make(chan error)
	r.lnd.LookupPaymentChannel <- LookupPaymentMessage{
		Hash:    hash,
		Updates: statusChan,
		Errors:  errorChan,
	}

	return statusChan, errorChan, nil
}
