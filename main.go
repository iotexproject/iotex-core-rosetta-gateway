// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
	"github.com/iotexproject/iotex-core-rosetta-gateway/services"
)

const (
	ConfigPath = "ConfigPath"
)

// NewBlockchainRouter returns a Mux http.Handler from a collection of
// Rosetta service controllers.
func NewBlockchainRouter(client ic.IoTexClient) (http.Handler, error) {
	asserter, err := asserter.NewServer(services.SupportedOperationTypes(),
		false,
		[]*types.NetworkIdentifier{
			&types.NetworkIdentifier{
				Blockchain: client.GetConfig().NetworkIdentifier.Blockchain,
				Network:    client.GetConfig().NetworkIdentifier.Network,
			},
		},
		[]string{},
		false,
		"",
	)
	if err != nil {
		return nil, err
	}
	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client), asserter)
	accountAPIController := server.NewAccountAPIController(services.NewAccountAPIService(client), asserter)
	blockAPIController := server.NewBlockAPIController(services.NewBlockAPIService(client), asserter)
	constructionAPIController := server.NewConstructionAPIController(services.NewConstructionAPIService(client), asserter)
	mempoolAPIController := server.NewMempoolAPIController(services.NewMemPoolAPIService(client), asserter)
	r := server.NewRouter(networkAPIController, accountAPIController, blockAPIController, constructionAPIController, mempoolAPIController)
	return server.CorsMiddleware(server.LoggerMiddleware(r)), nil
}

func main() {
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("ERROR: Failed to parse config: %v\n", err)
	}
	// Prepare a new gRPC client.
	client, err := ic.NewIoTexClient(cfg)
	if err != nil {
		log.Fatalf("ERROR: Failed to prepare IoTex gRPC client: %v\n", err)
	}

	// Start the server.
	router, err := NewBlockchainRouter(client)
	if err != nil {
		log.Fatalf("ERROR: Failed to init router: %v\n", err)
	}
	log.Println("listen", "0.0.0.0:"+cfg.Server.Port)
	if err := http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, router); err != nil {
		log.Fatalf("IoTex Rosetta Gateway server exited with error: %v\n", err)
	}
}
