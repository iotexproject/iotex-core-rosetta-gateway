// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex_client

import (
	"context"
	"encoding/hex"

	"github.com/coinbase/rosetta-sdk-go/types"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	// NonceKey is the name of the key in the Metadata map inside a
	// ConstructionMetadataResponse that specifies the next valid nonce.
	NonceKey = "nonce"
)

func fillIndex(transactions []*types.Transaction) []*types.Transaction {
	for i, t := range transactions {
		if len(t.Operations) == 0 {
			transactions = append(transactions[:i], transactions[i+1:]...)
			continue
		}
		for j, oper := range t.Operations {
			oper.OperationIdentifier.Index = int64(j)
		}
	}
	return transactions
}

func getTransactionLog(ctx context.Context, height int64, client iotexapi.APIServiceClient) (
	transferLogMap map[string][]*iotextypes.TransactionLog_Transaction, err error) {
	transferLogMap = make(map[string][]*iotextypes.TransactionLog_Transaction)
	transferLog, err := client.GetTransactionLogByBlockHeight(
		ctx,
		&iotexapi.GetTransactionLogByBlockHeightRequest{BlockHeight: uint64(height)},
	)
	if err != nil {
		return nil, err
	}

	for _, a := range transferLog.GetTransactionLogs().GetLogs() {
		h := hex.EncodeToString(a.ActionHash)
		transferLogMap[h] = a.GetTransactions()
	}
	return transferLogMap, nil
}

func getCaller(act *iotextypes.Action) (callerAddr address.Address, err error) {
	srcPub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		return
	}
	callerAddr, err = address.FromBytes(srcPub.Hash())
	return
}
