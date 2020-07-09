// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package inject

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

const (
	sender      = "io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02"
	privateKey  = "414efa99dfac6f4095d6954713fb0085268d400d6a05a8ae8a69b5b1c10b4bed"
	sender2     = "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms"
	privateKey2 = "cfa6ef757dee2e50351620dca002d32b9c090cfda55fb81f37f1d26b273743f1"
	to          = "io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg"
	receipt     = "io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms"
	endpoint    = "127.0.0.1:14014"
)

var (
	gasPrice = big.NewInt(0).SetUint64(1e12)
	gasLimit = uint64(10000000)
)

func TestInjectTransfer(t *testing.T) {
	for i := 0; i < 21; i++ {
		fmt.Println("inject transfer", i)
		injectTransfer(t)
	}
}

func TestMultisend(t *testing.T) {
	for i := 0; i < 21; i++ {
		fmt.Println("inject multisend contract", i)
		injectMultisend(t)
	}
}

func TestCandidateRegister(t *testing.T) {
	fmt.Println("inject candidate register")
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(privateKey2)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	getacc, err := c.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: sender2})
	require.NoError(err)
	fmt.Println("nonce:", getacc.AccountMeta.PendingNonce)
	cr, err := action.NewCandidateRegister(getacc.AccountMeta.PendingNonce, "xxxx", sender2, sender2, sender2,
		"12001000000000000000000000",
		7, false, nil, gasLimit, gasPrice)
	require.NoError(err)
	sk, err := crypto.HexStringToPrivateKey(privateKey2)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(getacc.AccountMeta.PendingNonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(cr).Build()
	selp, err := action.Sign(elp, sk)
	require.NoError(err)
	request := &iotexapi.SendActionRequest{Action: selp.Proto()}

	resp, err := c.API().SendAction(context.Background(), request)
	require.NoError(err)
	require.NotEmpty(resp.GetActionHash())
	checkHash(resp.GetActionHash(), t)
}

func TestStakeCreate(t *testing.T) {
	fmt.Println("inject stake create")
	stakeCreate(t, privateKey, sender, false)
	stakeCreate(t, privateKey2, sender2, true)
}

func stakeCreate(t *testing.T, pri, addr string, autostake bool) {
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(pri)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	getacc, err := c.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: addr})
	require.NoError(err)
	fmt.Println("nonce:", getacc.AccountMeta.PendingNonce)
	cr, err := action.NewCreateStake(getacc.AccountMeta.PendingNonce, "xxxx", "1200100000000000000000000", 0, autostake, nil, gasLimit, gasPrice)
	require.NoError(err)
	sk, err := crypto.HexStringToPrivateKey(pri)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(getacc.AccountMeta.PendingNonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(cr).Build()
	selp, err := action.Sign(elp, sk)
	require.NoError(err)
	request := &iotexapi.SendActionRequest{Action: selp.Proto()}

	resp, err := c.API().SendAction(context.Background(), request)
	require.NoError(err)
	require.NotEmpty(resp.GetActionHash())
	checkHash(resp.GetActionHash(), t)
}

func TestStakeAddDeposit(t *testing.T) {
	fmt.Println("inject stake add deposit")
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(privateKey2)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	getacc, err := c.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: sender2})
	require.NoError(err)
	cr, err := action.NewDepositToStake(getacc.AccountMeta.PendingNonce, 2, "1200100000000000000000000", nil,
		gasLimit, gasPrice)
	require.NoError(err)
	sk, err := crypto.HexStringToPrivateKey(privateKey2)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(getacc.AccountMeta.PendingNonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(cr).Build()
	selp, err := action.Sign(elp, sk)
	require.NoError(err)
	request := &iotexapi.SendActionRequest{Action: selp.Proto()}

	resp, err := c.API().SendAction(context.Background(), request)
	require.NoError(err)
	require.NotEmpty(resp.GetActionHash())
	checkHash(resp.GetActionHash(), t)
}

func TestStakeUnstake(t *testing.T) {
	fmt.Println("inject stake unstake")
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	getacc, err := c.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: sender})
	require.NoError(err)
	fmt.Println("nonce:", getacc.AccountMeta.PendingNonce)
	cr, err := action.NewUnstake(getacc.AccountMeta.PendingNonce, 1, nil,
		gasLimit, gasPrice)
	require.NoError(err)
	sk, err := crypto.HexStringToPrivateKey(privateKey)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(getacc.AccountMeta.PendingNonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(cr).Build()
	selp, err := action.Sign(elp, sk)
	require.NoError(err)
	request := &iotexapi.SendActionRequest{Action: selp.Proto()}

	resp, err := c.API().SendAction(context.Background(), request)
	require.NoError(err)
	require.NotEmpty(resp.GetActionHash())
	checkHash(resp.GetActionHash(), t)
}

func TestStakeWithdraw(t *testing.T) {
	fmt.Println("inject stake withdraw")
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	getacc, err := c.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: sender})
	require.NoError(err)
	fmt.Println("nonce:", getacc.AccountMeta.PendingNonce)
	cr, err := action.NewWithdrawStake(getacc.AccountMeta.PendingNonce, 1, nil, gasLimit, gasPrice)
	require.NoError(err)
	sk, err := crypto.HexStringToPrivateKey(privateKey)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(getacc.AccountMeta.PendingNonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(cr).Build()
	selp, err := action.Sign(elp, sk)
	require.NoError(err)
	request := &iotexapi.SendActionRequest{Action: selp.Proto()}

	resp, err := c.API().SendAction(context.Background(), request)
	require.NoError(err)
	require.NotEmpty(resp.GetActionHash())
	checkHash(resp.GetActionHash(), t)
}

func TestInjectTransferUseExecution(t *testing.T) {
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	getacc, err := c.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: sender})
	require.NoError(err)

	execution, err := action.NewExecution(to, getacc.AccountMeta.PendingNonce, big.NewInt(111), gasLimit, gasPrice, nil)
	require.NoError(err)
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(getacc.AccountMeta.PendingNonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(execution).Build()
	selp, err := action.Sign(elp, acc.PrivateKey())
	require.NoError(err)
	ret, err := iotexapi.NewAPIServiceClient(conn).SendAction(context.Background(), &iotexapi.SendActionRequest{
		Action: selp.Proto(),
	})
	require.NoError(err)
	checkHash(ret.ActionHash, t)
}

func TestGetImplicitLog(t *testing.T) {
	InContractTransfer := common.Hash{}
	BucketWithdrawAmount := hash.BytesToHash256([]byte("withdrawAmount"))
	fmt.Println("TestGetImplicitLog")
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	for i := uint64(1); i < 120; i++ {
		ret, err := c.API().GetImplicitTransferLogByBlockHeight(context.Background(),
			&iotexapi.GetImplicitTransferLogByBlockHeightRequest{
				BlockHeight: i})
		if err != nil {
			continue
		}
		for _, trans := range ret.GetBlockImplicitTransferLog().GetImplicitTransferLog() {
			for _, t := range trans.GetTransactions() {
				switch {
				case bytes.Compare(t.GetTopic(), InContractTransfer[:]) == 0:
					fmt.Println(i, "execution", t.Sender, t.Recipient, t.Amount)
				case bytes.Compare(t.GetTopic(), BucketWithdrawAmount[:]) == 0:
					fmt.Println(i, "stakewithdraw", t.Sender, t.Recipient, t.Amount)
				default:
					fmt.Println(i, "other")
				}
			}
		}
	}
}

func injectMultisend(t *testing.T) {
	require := require.New(t)
	contract := deployContract(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	abi, err := abi.JSON(strings.NewReader(MultisendABI))
	require.NoError(err)
	contractAddr, err := address.FromString(contract)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	r1, err := address.FromString(to)
	require.NoError(err)
	r2, err := address.FromString(receipt)
	require.NoError(err)
	r1ethAddress := common.HexToAddress(hex.EncodeToString(r1.Bytes()))
	r2ethAddress := common.HexToAddress(hex.EncodeToString(r2.Bytes()))
	hash, err := c.Contract(contractAddr, abi).Execute("multiSend", []common.Address{r1ethAddress, r2ethAddress}, []*big.Int{big.NewInt(1), big.NewInt(2)}, "").SetGasPrice(gasPrice).SetGasLimit(gasLimit).SetAmount(big.NewInt(3)).Call(context.Background())
	require.NoError(err)
	require.NotNil(hash)
	checkHash(hex.EncodeToString(hash[:]), t)
}

func injectTransfer(t *testing.T) {
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	to, err := address.FromString(to)
	require.NoError(err)
	hash, err := c.Transfer(to, big.NewInt(0).SetUint64(1000)).SetGasPrice(gasPrice).SetGasLimit(gasLimit).Call(context.Background())
	require.NoError(err)
	require.NotNil(hash)
	checkHash(hex.EncodeToString(hash[:]), t)
}

func deployContract(t *testing.T) string {
	require := require.New(t)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(privateKey)
	require.NoError(err)
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)

	data, err := hex.DecodeString(MultisendBin[2:])
	require.NoError(err)

	hash, err := c.DeployContract(data).SetGasPrice(gasPrice).SetGasLimit(gasLimit).Call(context.Background())
	require.NoError(err)
	require.NotNil(hash)
	fmt.Println("hash", hex.EncodeToString(hash[:]))
	time.Sleep(5 * time.Second)
	receiptResponse, err := c.GetReceipt(hash).Call(context.Background())
	require.NoError(err)
	contractAddress := receiptResponse.GetReceiptInfo().GetReceipt().GetContractAddress()
	fmt.Println("Status:", receiptResponse.GetReceiptInfo().GetReceipt().Status)
	fmt.Println("Contract Address:", contractAddress)
	return contractAddress
}

func checkHash(h string, t *testing.T) {
	fmt.Println("check hash:", h)
	require := require.New(t)
	time.Sleep(5 * time.Second)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(err)
	defer conn.Close()
	ha, err := hash.HexStringToHash256(h)
	require.NoError(err)
	c := iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	receiptResponse, err := c.GetReceipt(ha).Call(context.Background())
	s := receiptResponse.GetReceiptInfo().GetReceipt().GetStatus()
	fmt.Println("status:", s)
}
