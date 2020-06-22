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
			Index: s.client.GetConfig().Genesis_block_identifier.Index,
			Hash:  s.client.GetConfig().Genesis_block_identifier.Hash,
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
					Status:     ic.StatusSuccess,
					Successful: true,
				},
				{
					Status:     ic.StatusFail,
					Successful: false,
				},
			},
			OperationTypes: []string{
				ic.ActionTypeFee,
				ic.Transfer,
				ic.Execution,
				ic.DepositToRewardingFund,
				ic.ClaimFromRewardingFund,
				ic.StakeCreate,
				ic.StakeWithdraw,
				ic.StakeAddDeposit,
				ic.CandidateRegister},
			Errors: ErrorList,
		},
	}, nil
}
