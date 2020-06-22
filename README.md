# IoTeX Gateway for Rosetta

This repository implements the [Rosetta](https://github.com/coinbase/rosetta-sdk-go) for the [IoTeX](https://iotex.io) blockchain.

To build the server:

	make

To run tests:

	make test

To clean-up:

	make clean


`make test` will automatically download and build the [rosetta-cli](https://github.com/coinbase/rosetta-cli) ,then run
 the gateway and validate it using `rosetta-cli`.
