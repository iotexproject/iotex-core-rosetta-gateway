# IoTeX Gateway for Rosetta

This repository implements the [Rosetta](https://github.com/coinbase/rosetta-sdk-go) for the [IoTeX](https://iotex.io) blockchain.


Support Verions
| iotex-core | rosetta-specifications | rosetta-cli |
|----------|-------------|-------------|
| v1.1.0|v1.3.1|v0.4.1|


To build the server:

	make

To run tests:

	make test

`make test` will automatically download and build the [rosetta-cli](https://github.com/coinbase/rosetta-cli) ,then run the gateway and validate it using `rosetta-cli`. More test details can be found here: [tests](https://github.com/iotexproject/iotex-core-rosetta-gateway/tree/master/tests)

To clean-up:

	make clean

## Develop iotex-core-rosetta-gateway with Docker

To build the Docker image from your local repo:

	docker build -f ./docker/dev/Dockerfile . -t iotexproject/iotex-core-rosetta-gateway

To run the Docker image:

	docker run -p 8080:8080 -e "ConfigPath=/etc/iotex/config.yaml" iotexproject/iotex-core-rosetta-gateway
	

## Run IoTeX mainnet node and Rosetta Gateway in a Docker

Please refer to [Deployment](https://github.com/iotexproject/iotex-core-rosetta-gateway/blob/master/docker/deploy) here.
