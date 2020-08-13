// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
)

const (
	latestVersion = "v1.1.0"
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
			Blockchain: s.client.GetConfig().NetworkIdentifier.Blockchain,
			Network:    s.client.GetConfig().NetworkIdentifier.Network,
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

	blk, err := s.client.GetLatestBlock(ctx)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}
	genesisblk, err := s.client.GetBlock(ctx, 1)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}
	resp := &types.NetworkStatusResponse{
		CurrentBlockIdentifier: blk.BlockIdentifier,
		CurrentBlockTimestamp:  blk.Timestamp, // ms
		GenesisBlockIdentifier: genesisblk.BlockIdentifier,
		Peers:                  nil,
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
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}
	packageVersion := version.GetServerMeta().GetPackageVersion()
	if packageVersion == "" {
		packageVersion = latestVersion
	}
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: s.client.GetConfig().Server.RosettaVersion,
			NodeVersion:    packageVersion,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     ic.StatusSuccess,
					Successful: true,
				},
				{
					Status:     ic.StatusFail,
					Successful: false,
				},
			},
			OperationTypes: SupportedOperationTypes(),
			Errors:         ErrorList,
		},
	}, nil
}
