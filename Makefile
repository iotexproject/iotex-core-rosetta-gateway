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
ROSETTA_CLI_RELEASE=0.2.4
IOTEX_SERVER=iotex-core
IOTEX_SERVER_RELEASE=master
MAGENTA = ""
OFF = ""

HAVE_WGET := $(shell which wget > /dev/null && echo 1)
ifdef HAVE_WGET
    DOWNLOAD := wget --quiet --show-progress --progress=bar:force:noscroll -O
else
    HAVE_CURL := $(shell which curl > /dev/null && echo 1)
    ifdef HAVE_CURL
        DOWNLOAD := curl --progress-bar --location -o
    else
        $(error Please install wget or curl)
    endif
endif

default: build
all: clean build test

.PHONY: build
build:
	$(GOBUILD) -o ./$(BUILD_TARGET_SERVER) .

.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

tests/rosetta-cli.tar.gz:
	@echo "$(MAGENTA)*** Downloading rosetta-cli release $(ROSETTA_CLI_RELEASE)...$(OFF)\n"
	@$(DOWNLOAD) $@ https://github.com/coinbase/rosetta-cli/archive/v$(ROSETTA_CLI_RELEASE).tar.gz

tests/iotex-core.tar.gz:
	@echo "$(MAGENTA)*** Downloading iotex-core release $(IOTEX_SERVER_RELEASE)...$(OFF)\n"
	@$(DOWNLOAD) $@ https://github.com/iotexproject/iotex-core/archive/$(IOTEX_SERVER_RELEASE).tar.gz

tests/rosetta-cli: tests/rosetta-cli.tar.gz
	@echo "$(MAGENTA)*** Building rosetta-cli...$(OFF)\n"
	@tar -xf $< -C tests
	@cd tests/rosetta-cli-$(ROSETTA_CLI_RELEASE) && go build
	@cp tests/rosetta-cli-$(ROSETTA_CLI_RELEASE)/rosetta-cli tests

tests/server: tests/iotex-core.tar.gz
	@echo "$(MAGENTA)*** Building iotex-core...$(OFF)\n"
	@tar -xf $< -C tests
	@cd tests/iotex-core-$(IOTEX_SERVER_RELEASE) && make build
	@cp tests/iotex-core-$(IOTEX_SERVER_RELEASE)/bin/server tests

.PHONY: test
test: build tests/rosetta-cli tests/server
	@echo "Running tests...\n"
	@chmod +x ./tests/test.sh
	@cd tests&&./test.sh&&chmod -x ./test.sh

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf ./tests/rosetta* ./tests/iotex-core* ./tests/*.db ./tests/server ./tests/*.tar.gz
	@rm -rf ./$(BUILD_TARGET_SERVER)
	@rm -rf $(COV_REPORT) $(COV_HTML) $(LINT_LOG)
	@find . -name $(COV_OUT) -delete
	@find . -name $(TESTBED_COV_OUT) -delete
	@$(GOCLEAN) -i $(PKGS)