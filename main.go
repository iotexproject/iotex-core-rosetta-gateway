// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/coinbase/rosetta-sdk-go/server"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
	"github.com/iotexproject/iotex-core-rosetta-gateway/services"
)

const (
	ConfigPath = "ConfigPath"
)

// NewBlockchainRouter returns a Mux http.Handler from a collection of
// Rosetta service controllers.
func NewBlockchainRouter(client ic.IoTexClient) http.Handler {
	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client), nil)
	accountAPIController := server.NewAccountAPIController(services.NewAccountAPIService(client), nil)
	blockAPIController := server.NewBlockAPIController(services.NewBlockAPIService(client), nil)
	constructionAPIController := server.NewConstructionAPIController(services.NewConstructionAPIService(client), nil)
	return server.NewRouter(networkAPIController, accountAPIController, blockAPIController, constructionAPIController)
}

func main() {
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := config.New(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse config: %v\n", err)
		os.Exit(1)
	}
	// Prepare a new gRPC client.
	client, err := ic.NewIoTexClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to prepare IoTex gRPC client: %v\n", err)
		os.Exit(1)
	}

	// Start the server.
	router := NewBlockchainRouter(client)
	fmt.Println("listen", "0.0.0.0:"+cfg.Server.Port)
	err = http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "IoTex Rosetta Gateway server exited with error: %v\n", err)
		os.Exit(1)
	}
}
