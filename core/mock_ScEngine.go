// Code generated by mockery v1.0.0. DO NOT EDIT.

package core

import mock "github.com/stretchr/testify/mock"

// MockScEngine is an autogenerated mock type for the ScEngine type
type MockScEngine struct {
	mock.Mock
}

// DestroyEngine provides a mock function with given fields:
func (_m *MockScEngine) DestroyEngine() {
	_m.Called()
}

// Execute provides a mock function with given fields: function, args
func (_m *MockScEngine) Execute(function string, args string) string {
	ret := _m.Called(function, args)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(function, args)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetGeneratedTXs provides a mock function with given fields:
func (_m *MockScEngine) GetGeneratedTXs() []*Transaction {
	ret := _m.Called()

	var r0 []*Transaction
	if rf, ok := ret.Get(0).(func() []*Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Transaction)
		}
	}

	return r0
}

// ImportContractAddr provides a mock function with given fields: contractAddr
func (_m *MockScEngine) ImportContractAddr(contractAddr Address) {
	_m.Called(contractAddr)
}

// ImportCurrBlockHeight provides a mock function with given fields: currBlkHeight
func (_m *MockScEngine) ImportCurrBlockHeight(currBlkHeight uint64) {
	_m.Called(currBlkHeight)
}

// ImportLocalStorage provides a mock function with given fields: state
func (_m *MockScEngine) ImportLocalStorage(state *ScState) {
	_m.Called(state)
}

// ImportNodeAddress provides a mock function with given fields: addr
func (_m *MockScEngine) ImportNodeAddress(addr Address) {
	_m.Called(addr)
}

// ImportPrevUtxos provides a mock function with given fields: utxos
func (_m *MockScEngine) ImportPrevUtxos(utxos []*UTXO) {
	_m.Called(utxos)
}

// ImportRewardStorage provides a mock function with given fields: rewards
func (_m *MockScEngine) ImportRewardStorage(rewards map[string]string) {
	_m.Called(rewards)
}

// ImportSeed provides a mock function with given fields: seed
func (_m *MockScEngine) ImportSeed(seed int64) {
	_m.Called(seed)
}

// ImportSourceCode provides a mock function with given fields: source
func (_m *MockScEngine) ImportSourceCode(source string) {
	_m.Called(source)
}

// ImportSourceTXID provides a mock function with given fields: txid
func (_m *MockScEngine) ImportSourceTXID(txid []byte) {
	_m.Called(txid)
}

// ImportTransaction provides a mock function with given fields: tx
func (_m *MockScEngine) ImportTransaction(tx *Transaction) {
	_m.Called(tx)
}

func (_m *MockScEngine) ImportContractCreateUTXO(utxo *UTXO) {
	_m.Called(utxo)
}

// ImportUTXOs provides a mock function with given fields: utxos
func (_m *MockScEngine) ImportUTXOs(utxos []*UTXO) {
	_m.Called(utxos)
}
