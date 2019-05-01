module github.com/lightninglabs/loop

require (
	github.com/btcsuite/btcd v0.0.0-20190426011420-63f50db2f70a
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/coreos/bbolt v1.3.2
	github.com/fortytw2/leaktest v1.3.0
	github.com/golang/protobuf v1.3.1
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/jessevdk/go-flags v1.4.0
	github.com/lightningnetwork/lnd v0.5.1-beta.0.20190322220528-6ad8be25e1aa
	github.com/lightningnetwork/lnd/queue v1.0.1
	github.com/urfave/cli v1.20.0
	golang.org/x/net v0.0.0-20190313220215-9f648a60d977
	google.golang.org/genproto v0.0.0-20190307195333-5fe7a883aa19
	google.golang.org/grpc v1.19.0
	gopkg.in/macaroon.v2 v2.1.0
)

replace github.com/lightningnetwork/lnd => github.com/joostjager/lnd v0.4.1-beta.0.20190501134916-d218a62fb3a1
