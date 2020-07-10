########################################################################################################################
# Copyright (c) 2020 IoTeX Foundation
# This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
# warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
# permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
# License 2.0 that can be found in the LICENSE file.
########################################################################################################################

# Go parameters
GOCMD=go
GOLINT=golint
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_TARGET_SERVER=iotex-core-rosetta-gateway

default: build
all: clean build test

.PHONY: build
build:
	$(GOBUILD) -o ./$(BUILD_TARGET_SERVER) .

.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

.PHONY: test
test:
	@docker build -f ./docker/test/Dockerfile . -t iotexproject/iotex-core-rosetta-test
	@docker run --rm iotexproject/iotex-core-rosetta-test
	@docker rmi iotexproject/iotex-core-rosetta-test

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf ./tests/rosetta* ./tests/iotex-core* ./tests/*.db ./tests/server ./tests/*.tar.gz ./tests/*.vlog ./tests/LOCK ./tests/MANIFEST
	@rm -rf ./$(BUILD_TARGET_SERVER)
	@rm -rf $(COV_REPORT) $(COV_HTML) $(LINT_LOG)
	@find . -name $(COV_OUT) -delete
	@find . -name $(TESTBED_COV_OUT) -delete
	@$(GOCLEAN) -i $(PKGS)