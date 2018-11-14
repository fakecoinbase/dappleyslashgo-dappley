package sc

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -ldappleyv8
#include <stdlib.h>
#include "v8/engine.h"
//blockchain
bool  Cgo_VerifyAddressFunc(const char *address);
int  Cgo_TransferFunc(void *handler, const char *to, const char *amount, const char *tip);
//storage
char* Cgo_StorageGetFunc(void *address, const char *key);
int   Cgo_StorageSetFunc(void *address, const char *key, const char *value);
int   Cgo_StorageDelFunc(void *address, const char *key);
int	  Cgo_RecordRewardFunc(void *handler, const char *address, const char *amount);
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/core"
)

var (
	v8once        = sync.Once{}
	v8EngineList  = make(map[uint64]*V8Engine)
	storagesMutex = sync.RWMutex{}
	currHandler   = uint64(100)
)

type V8Engine struct {
	source        string
	storage       map[string]string
	rewards       map[string]string
	contractAddr  core.Address
	contractUTXOs []*core.UTXO
	sourceTXID    []byte
	generatedTXs  []*core.Transaction
	handler       uint64
}

func InitializeV8Engine() {
	C.Initialize()
	C.InitializeBlockchain(
		(C.FuncVerifyAddress)(unsafe.Pointer(C.Cgo_VerifyAddressFunc)),
		(C.FuncTransfer)(unsafe.Pointer(C.Cgo_TransferFunc)),
	)
	C.InitializeStorage(
		(C.FuncStorageGet)(unsafe.Pointer(C.Cgo_StorageGetFunc)),
		(C.FuncStorageSet)(unsafe.Pointer(C.Cgo_StorageSetFunc)),
		(C.FuncStorageDel)(unsafe.Pointer(C.Cgo_StorageDelFunc)),
	)
	C.InitializeRewardDistributor((C.FuncRecordReward)(unsafe.Pointer(C.Cgo_RecordRewardFunc)))
}

//NewV8Engine generates a new V8Engine instance
func NewV8Engine() *V8Engine {
	v8once.Do(func() { InitializeV8Engine() })
	engine := &V8Engine{
		source:       "",
		storage:      make(map[string]string),
		contractAddr: core.NewAddress(""),
		handler:      currHandler,
	}
	currHandler++
	storagesMutex.Lock()
	defer storagesMutex.Unlock()
	v8EngineList[engine.handler] = engine
	return engine
}

func (sc *V8Engine) ImportSourceCode(source string) {
	sc.source = source
}

func (sc *V8Engine) ImportLocalStorage(storage map[string]string) {
	sc.storage = storage
}

// ImportContractAddr supplies the wallet address of the contract to the engine
func (sc *V8Engine) ImportContractAddr(contractAddr core.Address) {
	sc.contractAddr = contractAddr
}

// ImportUTXOs supplies the list of contract's UTXOs to the engine
func (sc *V8Engine) ImportUTXOs(utxos []*core.UTXO) {
	sc.contractUTXOs = utxos
}

// ImportSourceTXID supplies the id of the transaction which executes the contract
func (sc *V8Engine) ImportSourceTXID(txid []byte) {
	sc.sourceTXID = txid
}

// GetGeneratedTXs returns the transactions generated as a result of executing the contract
func (sc *V8Engine) GetGeneratedTXs() []*core.Transaction {
	return sc.generatedTXs
}

func (sc *V8Engine) ImportRewardStorage(rewards map[string]string) {
	sc.rewards = rewards
}

func (sc *V8Engine) Execute(function, args string) string {
	res := "\"\""
	status := "success"
	var result *C.char

	cSource := C.CString(sc.source)
	defer C.free(unsafe.Pointer(cSource))
	C.InitializeSmartContract(cSource)

	functionCallScript := prepareFuncCallScript(function, args)
	cFunction := C.CString(functionCallScript)
	defer C.free(unsafe.Pointer(cFunction))

	if C.executeV8Script(cFunction, C.uintptr_t(sc.handler), &result) > 0 {
		status = "failed"
	}

	if result != nil {
		res = C.GoString(result)
		C.free(unsafe.Pointer(result))
	}

	logger.WithFields(logger.Fields{
		"result": res,
		"status": status,
	}).Info("Smart Contract Execution Ends.")
	return res
}

func prepareFuncCallScript(function, args string) string {
	return fmt.Sprintf(
		`var instance = new _native_require();instance["%s"].apply(instance, [%s]);`,
		function,
		args,
	)
}

func getV8EngineByAddress(handler uint64) *V8Engine {
	storagesMutex.Lock()
	defer storagesMutex.Unlock()
	return v8EngineList[handler]
}
