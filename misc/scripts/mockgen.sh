#!/bin/bash

mkdir -p ./iotex-client/mock_client
mockgen -destination=./iotex-client/mock_client/mock_client.go \
  -source=./iotex-client/iotex-client.go \
  -package=mock_client \
  IoTexClient
