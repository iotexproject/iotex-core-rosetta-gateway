// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex_client

import (
	"bytes"
	"context"
	"encoding/hex"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

func assertAction(act *iotextypes.Action, operations operationList) operationList {
	oper := &operation{}
	oper.amount = "0"
	switch {
	case act.GetCore().GetTransfer() != nil:
		oper.actionType = Transfer
		oper.amount = act.GetCore().GetTransfer().GetAmount()
		oper.dst = act.GetCore().GetTransfer().GetRecipient()
	case act.GetCore().GetDepositToRewardingFund() != nil:
		oper.actionType = DepositToRewardingFund
		oper.amount = act.GetCore().GetDepositToRewardingFund().GetAmount()
		oper.dst = RewardingAddress
	case act.GetCore().GetClaimFromRewardingFund() != nil:
		oper.actionType = ClaimFromRewardingFund
		oper.amount = act.GetCore().GetClaimFromRewardingFund().GetAmount()
		oper.isPositive = true
		oper.dst = RewardingAddress
	}
	if oper.amount != "0" && oper.actionType != "" {
		operations = append(operations, oper)
	}

	return operations
}

func getContractAddress(ctx context.Context, h string, client iotexapi.APIServiceClient) (contractAddr string, err error) {
	// need to get contract address generated of this action hash
	responseReceipt, err := client.GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{ActionHash: h})
	if err != nil {
		return
	}
	contractAddr = responseReceipt.GetReceiptInfo().GetReceipt().GetContractAddress()
	return
}

func getCaller(act *iotextypes.Action) (callerAddr address.Address, err error) {
	srcPub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		return
	}
	callerAddr, err = address.FromBytes(srcPub.Hash())
	return
}

func getActionType(topic []byte) string {
	InContractTransfer := common.Hash{}
	BucketWithdrawAmount := hash.BytesToHash256([]byte("withdrawAmount"))
	BucketCreateAmount := hash.BytesToHash256([]byte("createAmount"))
	BucketDepositAmount := hash.BytesToHash256([]byte("depositAmount"))
	CandidateRegistrationFee := hash.BytesToHash256([]byte("registrationFee"))
	CandidateSelfStake := hash.BytesToHash256([]byte("selfStake"))
	switch {
	case bytes.Compare(topic, InContractTransfer[:]) == 0:
		return Execution
	case bytes.Compare(topic, BucketWithdrawAmount[:]) == 0:
		return StakeWithdraw
	case bytes.Compare(topic, BucketCreateAmount[:]) == 0:
		return StakeCreate
	case bytes.Compare(topic, BucketDepositAmount[:]) == 0:
		return StakeAddDeposit
	case bytes.Compare(topic, CandidateRegistrationFee[:]) == 0:
		return CandidateRegister
	case bytes.Compare(topic, CandidateSelfStake[:]) == 0:
		return CandidateRegister
	}
	return ""
}

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

func getImplicitTransferLog(ctx context.Context, height int64, client iotexapi.APIServiceClient) (
	transferLogMap map[string][]*iotextypes.ImplicitTransferLog_Transaction, err error) {
	transferLogMap = make(map[string][]*iotextypes.ImplicitTransferLog_Transaction)
	transferLog, err := client.GetImplicitTransferLogByBlockHeight(
		ctx,
		&iotexapi.GetImplicitTransferLogByBlockHeightRequest{BlockHeight: uint64(height)},
	)

	if err == nil && transferLog.GetBlockImplicitTransferLog().GetNumTransactions() != 0 {
		for _, a := range transferLog.GetBlockImplicitTransferLog().GetImplicitTransferLog() {
			h := hex.EncodeToString(a.ActionHash)
			transferLogMap[h] = a.GetTransactions()
		}
	}
	return transferLogMap, nil
}
