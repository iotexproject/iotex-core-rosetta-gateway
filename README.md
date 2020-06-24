# IoTeX Gateway for Rosetta

This repository implements the [Rosetta](https://github.com/coinbase/rosetta-sdk-go) for the [IoTeX](https://iotex.io) blockchain.

To build the server:

	make

To run tests:

	make test

To clean-up:

	make clean


`make test` will automatically download and build the [rosetta-cli](https://github.com/coinbase/rosetta-cli) ,then run the gateway and validate it using `rosetta-cli`.

# Running iotex-core-rosetta-gateway in Docker

To build the Docker image:

	docker build -f ./docker/Dockerfile . -t iotexproject/iotex-core-rosetta-gateway

To run the Docker image:

	docker run -p 8080:8080 -e "ConfigPath=/etc/iotex/config.yaml" iotexproject/iotex-core-rosetta-gateway