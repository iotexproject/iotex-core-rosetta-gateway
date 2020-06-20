package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/coinbase/rosetta-sdk-go/server"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
	"github.com/iotexproject/iotex-core-rosetta-gateway/services"
)

const (
	// GatewayPort is the name of the environment variable that specifies
	// which port the IoTex Rosetta gateway should run on.
	GatewayPort = "GatewayPort"
	// IoTexChainPoint is the name of the environment variable that specifies
	// which the IoTex blockchain endpoint.
	IoTexChainPoint = "IoTexChainPoint"
	ConfigPath      = "ConfigPath"
)

// NewBlockchainRouter returns a Mux http.Handler from a collection of
// Rosetta service controllers.
func NewBlockchainRouter(client ic.IoTexClient) http.Handler {
	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client))
	accountAPIController := server.NewAccountAPIController(services.NewAccountAPIService(client))
	blockAPIController := server.NewBlockAPIController(services.NewBlockAPIService(client))
	constructionAPIController := server.NewConstructionAPIController(services.NewConstructionAPIService(client))
	return server.NewRouter(networkAPIController, accountAPIController, blockAPIController, constructionAPIController)
}

func main() {
	// Get server port from environment variable or use the default.
	port := os.Getenv(GatewayPort)
	if port == "" {
		port = "8080"
	}
	addr := os.Getenv(IoTexChainPoint)
	if addr == "" {
		fmt.Fprintf(os.Stderr, "ERROR: %s environment variable missing\n", IoTexChainPoint)
		os.Exit(1)
	}
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "config.json"
	}

	// Prepare a new gRPC client.
	client, err := ic.NewIoTexClient(addr, configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to prepare IoTex gRPC client: %v\n", err)
		os.Exit(1)
	}

	// Start the server.
	router := NewBlockchainRouter(client)
	fmt.Println("listen", "0.0.0.0:"+port)
	err = http.ListenAndServe("0.0.0.0:"+port, router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "IoTex Rosetta Gateway server exited with error: %v\n", err)
		os.Exit(1)
	}
}
