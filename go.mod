module github.com/iotexproject/iotex-core-rosetta-gateway

go 1.14

require (
	github.com/coinbase/rosetta-sdk-go v0.3.4
	github.com/elastic/gosigar v0.11.0 // indirect
	github.com/ethereum/go-ethereum v1.9.18
	github.com/golang/protobuf v1.4.2
	github.com/iotexproject/go-pkgs v0.1.2-0.20200523040337-5f1d9ddaa8ee
	github.com/iotexproject/iotex-address v0.2.2
	github.com/iotexproject/iotex-core v1.1.0
	github.com/iotexproject/iotex-proto v0.4.1
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/config v1.3.1
	go.uber.org/multierr v1.5.0 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.3.0
