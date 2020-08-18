package services

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/golang/mock/gomock"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"

	ic "github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client"
	"github.com/iotexproject/iotex-core-rosetta-gateway/iotex-client/mock_client"
)

func TestConstructionAPIService_ConstructionCombine(t *testing.T) {
	var (
		cfg     = testConfig()
		account = &types.AccountIdentifier{
			Address: "test account address",
		}
		publicKey = &types.PublicKey{
			Bytes:     []byte("hello"),
			CurveType: types.Secp256k1,
		}
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		signature = &types.Signature{
			SigningPayload: &types.SigningPayload{
				Address: account.Address,
				Bytes:   []byte("blah"),
			},
			PublicKey:     publicKey,
			SignatureType: types.Ed25519,
			Bytes:         []byte("hellohellohellohellohellohellohellohellohellohellohellohellohello"),
		}
		unsignedTransaction = "0a2e0801106552280a0731303030303030120e726563697069656e7420616464721a0d74657374207472616e73666572120568656c6c6f1a4168656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f"

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.ConstructionCombineResponse{SignedTransaction: unsignedTransaction}
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	resp, typErr := clt.ConstructionCombine(context.Background(), &types.ConstructionCombineRequest{
		NetworkIdentifier:   networkIdentifier,
		UnsignedTransaction: unsignedTransaction,
		Signatures:          []*types.Signature{signature},
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestConstructionAPIService_ConstructionDerive(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		publicKey = &types.PublicKey{CurveType: types.Secp256k1}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.ConstructionDeriveResponse{Address: "io13rjq2c07mqhe8sdd7nf9a4vcmnyk9mn72hu94e"}
	)
	pubKey, err := hex.DecodeString("04403d3c0dbd3270ddfc248c3df1f9aafd60f1d8e7456961c9ef26292262cc68f0ea9690263bef9e197a38f06026814fc70912c2b98d2e90a68f8ddc5328180a01")
	require.NoError(err)
	publicKey.Bytes = pubKey
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	resp, typErr := clt.ConstructionDerive(context.Background(), &types.ConstructionDeriveRequest{
		NetworkIdentifier: networkIdentifier,
		PublicKey:         publicKey,
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestConstructionAPIService_ConstructionHash(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.TransactionIdentifierResponse{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: "b08ab7d0c20384f06832db61b53674d00d33e82f41e5669b81c852ea0cb89297",
			},
		}
		signedTransaction = "0a630801100a18aa9c012214313130303030303030303030303030303030303052430a16313031303030303030303030303030303030303030301229696f316a6830656b6d63637977666b6d6a3765387173757a7375706e6c6b337735333337686a6a6732124104755ce6d8903f6b3793bddb4ea5d3589d637de2d209ae0ea930815c82db564ee8cc448886f639e8a0c7e94e99a5c1335b583c0bc76ef30dd6a1038ed9da8daf331a411eba3664e68f048d206c537d855f5d1853b9e9a0cff27c0dbfc5b60a58c738117ae0f9e2647afafb0051d582c493fa04c25b87c618f1081cc2d2593f6a34ab3a00"
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	resp, typErr := clt.ConstructionHash(context.Background(), &types.ConstructionHashRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedTransaction,
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestConstructionAPIService_ConstructionMetadata(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		gasLimit        = uint64(10000000)
		suggestGasPrice = uint64(1)
		nonce           = uint64(1)
		block           = &types.BlockIdentifier{
			Hash:  "block1",
			Index: 1,
		}
		amount = &types.Amount{
			Value: "100",
			Currency: &types.Currency{
				Symbol:   "IOTX",
				Decimals: 18,
			},
		}
		coin = &types.Coin{
			CoinIdentifier: &types.CoinIdentifier{
				Identifier: "IOTX",
			},
			Amount: amount,
		}
		accountBalanceResp = &types.AccountBalanceResponse{
			BlockIdentifier: block,
			Balances:        []*types.Amount{amount},
			Coins:           []*types.Coin{coin},
			Metadata:        map[string]interface{}{ic.NonceKey: nonce},
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.ConstructionMetadataResponse{
			Metadata: map[string]interface{}{
				"gasLimit":  gasLimit,
				"gasPrice":  suggestGasPrice,
				ic.NonceKey: nonce,
			},
			SuggestedFee: []*types.Amount{
				{
					Value:    "10000000", // gasLimit * suggestGasPrice
					Currency: amount.Currency,
				},
			},
		}
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(0)), gomock.AssignableToTypeOf("")).
		Return(accountBalanceResp, nil).AnyTimes()
	cli.EXPECT().EstimateGasForAction(gomock.Any(), gomock.AssignableToTypeOf(&iotextypes.Action{})).
		Return(gasLimit, nil).AnyTimes()
	cli.EXPECT().SuggestGasPrice(gomock.Any()).Return(suggestGasPrice, nil).AnyTimes()
	resp, typErr := clt.ConstructionMetadata(context.Background(), &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options: map[string]interface{}{
			"sender": "io13rjq2c07mqhe8sdd7nf9a4vcmnyk9mn72hu94e",
			"type":   iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
		},
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestConstructionAPIService_ConstructionParse(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		currency = &types.Currency{
			Symbol:   "IOTX",
			Decimals: 18,
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		tests   = []struct {
			req *types.ConstructionParseRequest
			rep *types.ConstructionParseResponse
			err *types.Error
		}{
			{
				req: &types.ConstructionParseRequest{
					NetworkIdentifier: networkIdentifier,
					Signed:            false,
					Transaction:       "0a2a080110011880ade204220131521c0a033130301a1574657374207472616e73666572207061796c6f6164",
				},
				rep: &types.ConstructionParseResponse{
					Operations: []*types.Operation{
						{
							OperationIdentifier: &types.OperationIdentifier{Index: 0},
							Type:                iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
							Account:             &types.AccountIdentifier{},
							Amount: &types.Amount{
								Value:    "-100",
								Currency: currency,
							},
						}, {
							OperationIdentifier: &types.OperationIdentifier{Index: 1},
							RelatedOperations: []*types.OperationIdentifier{
								{Index: 0},
							},
							Type:    iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
							Account: &types.AccountIdentifier{},
							Amount: &types.Amount{
								Value:    "100",
								Currency: currency,
							},
						},
					},
					Metadata: map[string]interface{}{
						"gasLimit": uint64(10000000),
						"gasPrice": uint64(1),
						"nonce":    uint64(1),
					},
				},
				err: nil,
			},
			// signed case
			{
				req: &types.ConstructionParseRequest{
					NetworkIdentifier: networkIdentifier,
					Signed:            true,
					Transaction:       "0a630801100a18aa9c012214313130303030303030303030303030303030303052430a16313031303030303030303030303030303030303030301229696f316a6830656b6d63637977666b6d6a3765387173757a7375706e6c6b337735333337686a6a6732124104755ce6d8903f6b3793bddb4ea5d3589d637de2d209ae0ea930815c82db564ee8cc448886f639e8a0c7e94e99a5c1335b583c0bc76ef30dd6a1038ed9da8daf331a411eba3664e68f048d206c537d855f5d1853b9e9a0cff27c0dbfc5b60a58c738117ae0f9e2647afafb0051d582c493fa04c25b87c618f1081cc2d2593f6a34ab3a00",
				},
				rep: &types.ConstructionParseResponse{
					Operations: []*types.Operation{
						{
							OperationIdentifier: &types.OperationIdentifier{Index: 0},
							Type:                iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
							Account: &types.AccountIdentifier{
								Address: "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms",
							},
							Amount: &types.Amount{
								Value:    "-1010000000000000000000",
								Currency: currency,
							},
						}, {
							OperationIdentifier: &types.OperationIdentifier{Index: 1},
							RelatedOperations: []*types.OperationIdentifier{
								{Index: 0},
							},
							Type: iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
							Account: &types.AccountIdentifier{
								Address: "io1jh0ekmccywfkmj7e8qsuzsupnlk3w5337hjjg2",
							},
							Amount: &types.Amount{
								Value:    "1010000000000000000000",
								Currency: currency,
							},
							CoinChange: nil,
							Metadata:   nil,
						},
					},
					Signers: []string{"io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms"},
					Metadata: map[string]interface{}{
						"gasLimit": uint64(20010),
						"gasPrice": uint64(11000000000000000000),
						"nonce":    uint64(10),
					},
				},
				err: nil,
			},
		}
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	for i, test := range tests {
		resp, typErr := clt.ConstructionParse(context.Background(), test.req)
		require.Equal(test.err, typErr, "index: %d", i)
		require.Equal(test.rep, resp, "index: %d", i)
	}
}

func TestConstructionAPIService_ConstructionPayloads(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		currency = &types.Currency{
			Symbol:   "IOTX",
			Decimals: 18,
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.ConstructionPayloadsResponse{
			UnsignedTransaction: "0a61100a18aa9c012214313130303030303030303030303030303030303052430a16313031303030303030303030303030303030303030301229696f316a6830656b6d63637977666b6d6a3765387173757a7375706e6c6b337735333337686a6a67321229696f316d666c70396d366863676d327163676863687364716a337a33656363726e656b783970306d73",
			Payloads: []*types.SigningPayload{
				{
					Address: "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms",
					Bytes: []byte{0x06, 0x63, 0x95, 0xdb, 0xab, 0xbe, 0x8d, 0x30, 0x1e, 0x00, 0xf8, 0x5b, 0x8f, 0xa8, 0xb9, 0x3f,
						0x90, 0xfd, 0x33, 0x8a, 0x8d, 0x7a, 0x52, 0x43, 0x7f, 0xa1, 0x21, 0x3a, 0xa1, 0x41, 0x9d, 0x8b},
					SignatureType: "ecdsa_recovery",
				},
			},
		}
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	resp, typErr := clt.ConstructionPayloads(context.Background(), &types.ConstructionPayloadsRequest{
		NetworkIdentifier: networkIdentifier,
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{Index: 0},
				Type:                iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
				Account: &types.AccountIdentifier{
					Address: "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms",
				},
				Amount: &types.Amount{
					Value:    "-1010000000000000000000",
					Currency: currency,
				},
			}, {
				OperationIdentifier: &types.OperationIdentifier{Index: 1},
				RelatedOperations: []*types.OperationIdentifier{
					{Index: 0},
				},
				Type: iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
				Account: &types.AccountIdentifier{
					Address: "io1jh0ekmccywfkmj7e8qsuzsupnlk3w5337hjjg2",
				},
				Amount: &types.Amount{
					Value:    "1010000000000000000000",
					Currency: currency,
				},
				CoinChange: nil,
				Metadata:   nil,
			},
		},
		Metadata: map[string]interface{}{
			"gasLimit": uint64(20010),
			"gasPrice": uint64(11000000000000000000),
			"nonce":    uint64(10),
		},
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestConstructionAPIService_ConstructionPreprocess(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}
		currency = &types.Currency{
			Symbol:   "IOTX",
			Decimals: 18,
		}
		maxFee = &types.Amount{
			Value:    "1000",
			Currency: currency,
		}
		suggestedFee = 1.23

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.ConstructionPreprocessResponse{
			Options: map[string]interface{}{
				"amount":        "1010000000000000000000",
				"decimals":      int32(18),
				"feeMultiplier": float64(1.23),
				"gasLimit":      uint64(20010),
				"gasPrice":      uint64(11000000000000000000),
				"maxFee":        "1000",
				"recipient":     "io1jh0ekmccywfkmj7e8qsuzsupnlk3w5337hjjg2",
				"sender":        "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms",
				"symbol":        "IOTX",
				"type":          "NATIVE_TRANSFER",
			},
		}
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	resp, typErr := clt.ConstructionPreprocess(context.Background(), &types.ConstructionPreprocessRequest{
		NetworkIdentifier: networkIdentifier,
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{Index: 0},
				Type:                iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
				Account: &types.AccountIdentifier{
					Address: "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms",
				},
				Amount: &types.Amount{
					Value:    "-1010000000000000000000",
					Currency: currency,
				},
			}, {
				OperationIdentifier: &types.OperationIdentifier{Index: 1},
				RelatedOperations: []*types.OperationIdentifier{
					{Index: 0},
				},
				Type: iotextypes.TransactionLogType_NATIVE_TRANSFER.String(),
				Account: &types.AccountIdentifier{
					Address: "io1jh0ekmccywfkmj7e8qsuzsupnlk3w5337hjjg2",
				},
				Amount: &types.Amount{
					Value:    "1010000000000000000000",
					Currency: currency,
				},
			},
		},
		Metadata: map[string]interface{}{
			"gasLimit": uint64(20010),
			"gasPrice": uint64(11000000000000000000),
			"nonce":    uint64(10),
		},
		MaxFee:                 []*types.Amount{maxFee},
		SuggestedFeeMultiplier: &suggestedFee,
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}

func TestConstructionAPIService_ConstructionSubmit(t *testing.T) {
	var (
		cfg               = testConfig()
		networkIdentifier = &types.NetworkIdentifier{
			Blockchain: "IoTeX",
			Network:    "testnet",
		}

		require = require.New(t)
		ctrl    = gomock.NewController(t)
		cli     = mock_client.NewMockIoTexClient(ctrl)
		clt     = NewConstructionAPIService(cli)
		ret     = &types.TransactionIdentifierResponse{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: "tx id",
			},
		}
		signedTransaction = "0a2e0801106552280a0731303030303030120e726563697069656e7420616464721a0d74657374207472616e73666572120568656c6c6f1a4168656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f68656c6c6f"
	)
	cli.EXPECT().GetConfig().Return(cfg).AnyTimes()
	cli.EXPECT().SubmitTx(gomock.Any(), gomock.AssignableToTypeOf(&iotextypes.Action{})).
		Return("tx id", nil).AnyTimes()
	resp, typErr := clt.ConstructionSubmit(context.Background(), &types.ConstructionSubmitRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedTransaction,
	})
	require.Nil(typErr)
	require.Equal(ret, resp)
}
