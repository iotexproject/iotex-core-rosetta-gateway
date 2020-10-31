// Code generated by MockGen. DO NOT EDIT.
// Source: ./iotex-client/iotex-client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	types "github.com/coinbase/rosetta-sdk-go/types"
	gomock "github.com/golang/mock/gomock"
	config "github.com/iotexproject/iotex-core-rosetta-gateway/config"
	iotexapi "github.com/iotexproject/iotex-proto/golang/iotexapi"
	iotextypes "github.com/iotexproject/iotex-proto/golang/iotextypes"
	reflect "reflect"
)

// MockIoTexClient is a mock of IoTexClient interface
type MockIoTexClient struct {
	ctrl     *gomock.Controller
	recorder *MockIoTexClientMockRecorder
}

// MockIoTexClientMockRecorder is the mock recorder for MockIoTexClient
type MockIoTexClientMockRecorder struct {
	mock *MockIoTexClient
}

// NewMockIoTexClient creates a new mock instance
func NewMockIoTexClient(ctrl *gomock.Controller) *MockIoTexClient {
	mock := &MockIoTexClient{ctrl: ctrl}
	mock.recorder = &MockIoTexClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIoTexClient) EXPECT() *MockIoTexClientMockRecorder {
	return m.recorder
}

// GetChainID mocks base method
func (m *MockIoTexClient) GetChainID(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainID", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChainID indicates an expected call of GetChainID
func (mr *MockIoTexClientMockRecorder) GetChainID(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainID", reflect.TypeOf((*MockIoTexClient)(nil).GetChainID), ctx)
}

// GetBlock mocks base method
func (m *MockIoTexClient) GetBlock(ctx context.Context, height int64) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", ctx, height)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock
func (mr *MockIoTexClientMockRecorder) GetBlock(ctx, height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockIoTexClient)(nil).GetBlock), ctx, height)
}

// GetLatestBlock mocks base method
func (m *MockIoTexClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestBlock", ctx)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestBlock indicates an expected call of GetLatestBlock
func (mr *MockIoTexClientMockRecorder) GetLatestBlock(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestBlock", reflect.TypeOf((*MockIoTexClient)(nil).GetLatestBlock), ctx)
}

// GetGenesisBlock mocks base method
func (m *MockIoTexClient) GetGenesisBlock(ctx context.Context) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenesisBlock", ctx)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenesisBlock indicates an expected call of GetGenesisBlock
func (mr *MockIoTexClientMockRecorder) GetGenesisBlock(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenesisBlock", reflect.TypeOf((*MockIoTexClient)(nil).GetGenesisBlock), ctx)
}

// GetAccount mocks base method
func (m *MockIoTexClient) GetAccount(ctx context.Context, height int64, owner string) (*types.AccountBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, height, owner)
	ret0, _ := ret[0].(*types.AccountBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockIoTexClientMockRecorder) GetAccount(ctx, height, owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockIoTexClient)(nil).GetAccount), ctx, height, owner)
}

// SubmitTx mocks base method
func (m *MockIoTexClient) SubmitTx(ctx context.Context, tx *iotextypes.Action) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitTx", ctx, tx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitTx indicates an expected call of SubmitTx
func (mr *MockIoTexClientMockRecorder) SubmitTx(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTx", reflect.TypeOf((*MockIoTexClient)(nil).SubmitTx), ctx, tx)
}

// GetStatus mocks base method
func (m *MockIoTexClient) GetStatus(ctx context.Context) (*iotexapi.GetChainMetaResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", ctx)
	ret0, _ := ret[0].(*iotexapi.GetChainMetaResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus
func (mr *MockIoTexClientMockRecorder) GetStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockIoTexClient)(nil).GetStatus), ctx)
}

// GetVersion mocks base method
func (m *MockIoTexClient) GetVersion(ctx context.Context) (*iotexapi.GetServerMetaResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion", ctx)
	ret0, _ := ret[0].(*iotexapi.GetServerMetaResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersion indicates an expected call of GetVersion
func (mr *MockIoTexClientMockRecorder) GetVersion(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockIoTexClient)(nil).GetVersion), ctx)
}

// GetTransactions mocks base method
func (m *MockIoTexClient) GetTransactions(ctx context.Context, height int64) ([]*types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", ctx, height)
	ret0, _ := ret[0].([]*types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions
func (mr *MockIoTexClientMockRecorder) GetTransactions(ctx, height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockIoTexClient)(nil).GetTransactions), ctx, height)
}

// GetConfig mocks base method
func (m *MockIoTexClient) GetConfig() *config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(*config.Config)
	return ret0
}

// GetConfig indicates an expected call of GetConfig
func (mr *MockIoTexClientMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockIoTexClient)(nil).GetConfig))
}

// SuggestGasPrice mocks base method
func (m *MockIoTexClient) SuggestGasPrice(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestGasPrice", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestGasPrice indicates an expected call of SuggestGasPrice
func (mr *MockIoTexClientMockRecorder) SuggestGasPrice(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestGasPrice", reflect.TypeOf((*MockIoTexClient)(nil).SuggestGasPrice), ctx)
}

// EstimateGasForAction mocks base method
func (m *MockIoTexClient) EstimateGasForAction(ctx context.Context, action *iotextypes.Action) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateGasForAction", ctx, action)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EstimateGasForAction indicates an expected call of EstimateGasForAction
func (mr *MockIoTexClientMockRecorder) EstimateGasForAction(ctx, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateGasForAction", reflect.TypeOf((*MockIoTexClient)(nil).EstimateGasForAction), ctx, action)
}

// GetBlockTransaction mocks base method
func (m *MockIoTexClient) GetBlockTransaction(ctx context.Context, actionHash string) (*types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockTransaction", ctx, actionHash)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockTransaction indicates an expected call of GetBlockTransaction
func (mr *MockIoTexClientMockRecorder) GetBlockTransaction(ctx, actionHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockTransaction", reflect.TypeOf((*MockIoTexClient)(nil).GetBlockTransaction), ctx, actionHash)
}

// GetMemPool mocks base method
func (m *MockIoTexClient) GetMemPool(ctx context.Context, actionHashes []string) ([]*types.TransactionIdentifier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMemPool", ctx, actionHashes)
	ret0, _ := ret[0].([]*types.TransactionIdentifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMemPool indicates an expected call of GetMemPool
func (mr *MockIoTexClientMockRecorder) GetMemPool(ctx, actionHashes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMemPool", reflect.TypeOf((*MockIoTexClient)(nil).GetMemPool), ctx, actionHashes)
}
