#!/usr/bin/env bash
go build -o ./iotex-core-rosetta-gateway .
export IoTexChainPoint=api.testnet.iotex.one:443
./iotex-core-rosetta-gateway
