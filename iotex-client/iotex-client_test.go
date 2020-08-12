package iotex_client

import (
	"context"
	"net"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotexapi/mock_iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-core-rosetta-gateway/config"
)

const (
	serverAddr = ":14014"
)

var (
	testCfg = &config.Config{
		NetworkIdentifier: config.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		},
		Currency: config.Currency{
			Symbol:   "IOTX",
			Decimals: 18,
		},
		Server: config.Server{
			Port:           "8080",
			Endpoint:       "localhost" + serverAddr,
			SecureEndpoint: false,
			RosettaVersion: "1.3.5",
		},
		KeepNoneTxAction: false,
	}
)

func newMockAPIServiceServer(r *require.Assertions, ctrl *gomock.Controller) *mock_iotexapi.MockAPIServiceServer {
	service := mock_iotexapi.NewMockAPIServiceServer(ctrl)
	server := grpc.NewServer()
	iotexapi.RegisterAPIServiceServer(server, service)
	listener, err := net.Listen("tcp", serverAddr)
	r.NoError(err)
	go func() {
		err := server.Serve(listener)
		r.NoError(err)
	}()
	return service
}

func TestIoTexClient_GetBlock(t *testing.T) {
	var (
		blkMetas = []*iotextypes.BlockMeta{
			{
				Hash:      "hash 1",
				Height:    1,
				Timestamp: &timestamp.Timestamp{},
			}, {
				Hash:      "hash 2",
				Height:    2,
				Timestamp: &timestamp.Timestamp{},
			},
		}
		expectBlk = &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Index: 2,
				Hash:  "hash 2",
			},
			ParentBlockIdentifier: &types.BlockIdentifier{
				Index: 1,
				Hash:  "hash 1",
			},
		}
	)
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	server := newMockAPIServiceServer(require, ctrl)
	server.EXPECT().
		GetBlockMetas(gomock.Any(), gomock.AssignableToTypeOf(&iotexapi.GetBlockMetasRequest{})).
		Return(&iotexapi.GetBlockMetasResponse{BlkMetas: blkMetas}, nil).
		AnyTimes()

	cli, err := NewIoTexClient(testCfg)
	require.NoError(err)
	block, err := cli.GetBlock(context.Background(), 2)
	require.NoError(err)
	require.Equal(expectBlk, block)
}
