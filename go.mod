module github.com/iotexproject/iotex-core-rosetta-gateway

go 1.14

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/coinbase/rosetta-sdk-go v0.3.4
	github.com/elastic/gosigar v0.11.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/iotexproject/go-pkgs v0.1.2-0.20200523040337-5f1d9ddaa8ee
	github.com/iotexproject/iotex-address v0.2.2
	github.com/iotexproject/iotex-proto v0.4.2
	github.com/pkg/errors v0.8.1
	github.com/spf13/cast v1.3.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	go.uber.org/config v1.3.1
	go.uber.org/multierr v1.5.0 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.3.0
