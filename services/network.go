package services

import (
	"context"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

type networkAPIService struct {
	client ic.IoTexClient
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(client ic.IoTexClient) server.NetworkAPIServicer {
	return &networkAPIService{
		client: client,
	}
}

// NetworkList implements the /network/list endpoint.
func (s *networkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{{
			Blockchain: s.client.GetConfig().Network_identifier.Blockchain,
			Network:    s.client.GetConfig().Network_identifier.Network,
		},
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *networkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	status, err := s.client.GetStatus(ctx)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}
	hei := int64(status.GetChainMeta().GetHeight())
	blk, err := s.client.GetBlock(ctx, hei)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}

	resp := &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: hei,
			Hash:  blk.Hash,
		},
		CurrentBlockTimestamp: blk.Timestamp, // ms
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: 1,
			Hash:  "",
		},
		Peers: nil,
	}

	return resp, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *networkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, s.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	version, err := s.client.GetVersion(ctx)
	fmt.Println(version, err)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}

	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: "1.3.5",
			NodeVersion:    version.GetServerMeta().GetPackageVersion(),
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     "OK",
					Successful: true,
				},
			},
			OperationTypes: []string{"transfer"},
			Errors:         ErrorList,
		},
	}, nil
}
