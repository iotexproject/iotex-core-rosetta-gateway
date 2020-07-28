module github.com/iotexproject/iotex-core-rosetta-gateway

go 1.13

require (
	github.com/coinbase/rosetta-sdk-go v0.1.5
	github.com/ethereum/go-ethereum v1.9.15 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/iotexproject/go-pkgs v0.1.2-0.20200523040337-5f1d9ddaa8ee
	github.com/iotexproject/iotex-proto v0.3.2-0.20200727212950-88e68ce8b8a7
	github.com/pkg/errors v0.8.1
	go.uber.org/config v1.3.1
	go.uber.org/multierr v1.5.0 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.3.0
