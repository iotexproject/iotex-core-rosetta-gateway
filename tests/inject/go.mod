module github.com/iotexproject/iotex-core-rosetta-gateway

go 1.13

require (
	github.com/elastic/gosigar v0.11.0 // indirect
	github.com/ethereum/go-ethereum v1.8.27
	github.com/iotexproject/go-pkgs v0.1.2-0.20200523040337-5f1d9ddaa8ee
	github.com/iotexproject/iotex-address v0.2.1
	github.com/iotexproject/iotex-antenna-go/v2 v2.4.1-0.20200814192756-b7914940625f
	github.com/iotexproject/iotex-proto v0.4.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.27.0
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.3.0
