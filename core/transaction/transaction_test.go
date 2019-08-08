// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either pubKeyHash 3 of the License, or
// (at your option) any later pubKeyHash.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/core/transaction/pb"
	"github.com/dappley/go-dappley/core/transaction_base"
	"github.com/dappley/go-dappley/core/utxo"
	"github.com/dappley/go-dappley/logic/transaction_logic"
	"github.com/dappley/go-dappley/logic/utxo_logic"
	"testing"

	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/core/account"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
	"github.com/dappley/go-dappley/storage"
	"github.com/dappley/go-dappley/util"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getAoB(length int64) []byte {
	return util.GenerateRandomAoB(length)
}

func GenerateFakeTxInputs() []transaction_base.TXInput {
	return []transaction_base.TXInput{
		{getAoB(2), 10, getAoB(2), getAoB(2)},
		{getAoB(2), 5, getAoB(2), getAoB(2)},
	}
}

func GenerateFakeTxOutputs() []transaction_base.TXOutput {
	return []transaction_base.TXOutput{
		{common.NewAmount(1), account.PubKeyHash(getAoB(2)), ""},
		{common.NewAmount(2), account.PubKeyHash(getAoB(2)), ""},
	}
}

var tx1 = Transaction{
	ID:       util.GenerateRandomAoB(1),
	Vin:      GenerateFakeTxInputs(),
	Vout:     GenerateFakeTxOutputs(),
	Tip:      common.NewAmount(5),
	GasLimit: common.NewAmount(0),
	GasPrice: common.NewAmount(0),
}

func TestTrimmedCopy(t *testing.T) {
	var tx1 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(2),
	}

	t2 := tx1.TrimmedCopy(false)

	t3 := transaction_logic.NewCoinbaseTX(account.NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"), "", 0, common.NewAmount(0))
	t4 := t3.TrimmedCopy(false)
	assert.Equal(t, tx1.ID, t2.ID)
	assert.Equal(t, tx1.Tip, t2.Tip)
	assert.Equal(t, tx1.Vout, t2.Vout)
	for index, vin := range t2.Vin {
		assert.Nil(t, vin.Signature)
		assert.Nil(t, vin.PubKey)
		assert.Equal(t, tx1.Vin[index].Txid, vin.Txid)
		assert.Equal(t, tx1.Vin[index].Vout, vin.Vout)
	}

	assert.Equal(t, t3.ID, t4.ID)
	assert.Equal(t, t3.Tip, t4.Tip)
	assert.Equal(t, t3.Vout, t4.Vout)
	for index, vin := range t4.Vin {
		assert.Nil(t, vin.Signature)
		assert.Nil(t, vin.PubKey)
		assert.Equal(t, t3.Vin[index].Txid, vin.Txid)
		assert.Equal(t, t3.Vin[index].Vout, vin.Vout)
	}
}

func TestSign(t *testing.T) {
	// Fake a key pair
	privKey, _ := ecdsa.GenerateKey(secp256k1.S256(), bytes.NewReader([]byte("fakefakefakefakefakefakefakefakefakefake")))
	ecdsaPubKey, _ := secp256k1.FromECDSAPublicKey(&privKey.PublicKey)
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	pubKeyHash, _ := account.NewUserPubKeyHash(pubKey)

	// Previous transactions containing UTXO of the Address
	prevTXs := []*utxo.UTXO{
		{transaction_base.TXOutput{common.NewAmount(13), pubKeyHash, ""}, []byte("01"), 0, utxo.UtxoNormal},
		{transaction_base.TXOutput{common.NewAmount(13), pubKeyHash, ""}, []byte("02"), 0, utxo.UtxoNormal},
		{transaction_base.TXOutput{common.NewAmount(13), pubKeyHash, ""}, []byte("03"), 0, utxo.UtxoNormal},
	}

	// New transaction to be signed (paid from the fake account)
	txin := []transaction_base.TXInput{
		{[]byte{1}, 0, nil, pubKey},
		{[]byte{3}, 0, nil, pubKey},
		{[]byte{3}, 2, nil, pubKey},
	}
	txout := []transaction_base.TXOutput{
		{common.NewAmount(19), pubKeyHash, ""},
	}
	tx := Transaction{nil, txin, txout, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}

	// transaction_logic.Sign the transaction
	err := transaction_logic.Sign(*privKey, prevTXs, &tx)
	if assert.Nil(t, err) {
		// Assert that the signatures were created by the fake key pair
		for i, vin := range tx.Vin {

			if assert.NotNil(t, vin.Signature) {
				txCopy := tx.TrimmedCopy(false)
				txCopy.Vin[i].Signature = nil
				txCopy.Vin[i].PubKey = []byte(pubKeyHash)

				verified, err := secp256k1.Verify(txCopy.Hash(), vin.Signature, ecdsaPubKey)
				assert.Nil(t, err)
				assert.True(t, verified)
			}
		}
	}
}

func TestVerifyCoinbaseTransaction(t *testing.T) {
	var prevTXs = map[string]Transaction{}

	var tx1 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(2),
	}

	var tx2 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(5),
	}
	var tx3 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(10),
	}
	var tx4 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(20),
	}
	prevTXs[string(tx1.ID)] = tx2
	prevTXs[string(tx2.ID)] = tx3
	prevTXs[string(tx3.ID)] = tx4

	// test verifying coinbase transactions
	var t5 = transaction_logic.NewCoinbaseTX(account.NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"), "", 5, common.NewAmount(0))
	bh1 := make([]byte, 8)
	binary.BigEndian.PutUint64(bh1, 5)
	txin1 := transaction_base.TXInput{nil, -1, bh1, []byte("Reward to test")}
	txout1 := transaction_base.NewTXOutput(common.NewAmount(10000000), account.NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"))
	var t6 = Transaction{nil, []transaction_base.TXInput{txin1}, []transaction_base.TXOutput{*txout1}, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}

	// test valid coinbase transaction
	_, err5 := transaction_logic.VerifyTransaction(&utxo_logic.UTXOIndex{}, &t5, 5)
	assert.Nil(t, err5)
	_, err6 := transaction_logic.VerifyTransaction(&utxo_logic.UTXOIndex{}, &t6, 5)
	assert.Nil(t, err6)

	// test coinbase transaction with incorrect blockHeight
	_, err5 = transaction_logic.VerifyTransaction(&utxo_logic.UTXOIndex{}, &t5, 10)
	assert.NotNil(t, err5)

	// test coinbase transaction with incorrect Subsidy
	bh2 := make([]byte, 8)
	binary.BigEndian.PutUint64(bh2, 5)
	txin2 := transaction_base.TXInput{nil, -1, bh2, []byte(nil)}
	txout2 := transaction_base.NewTXOutput(common.NewAmount(9), account.NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"))
	var t7 = Transaction{nil, []transaction_base.TXInput{txin2}, []transaction_base.TXOutput{*txout2}, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}
	_, err7 := transaction_logic.VerifyTransaction(&utxo_logic.UTXOIndex{}, &t7, 5)
	assert.NotNil(t, err7)

}

func TestVerifyNoCoinbaseTransaction(t *testing.T) {
	// Fake a key pair
	privKey, _ := ecdsa.GenerateKey(secp256k1.S256(), bytes.NewReader([]byte("fakefakefakefakefakefakefakefakefakefake")))
	privKeyByte, _ := secp256k1.FromECDSAPrivateKey(privKey)
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	pubKeyHash, _ := account.NewUserPubKeyHash(pubKey)
	//Address := KeyPair{*privKey, pubKey}.GenerateAddress()

	// Fake a wrong key pair
	wrongPrivKey, _ := ecdsa.GenerateKey(secp256k1.S256(), bytes.NewReader([]byte("FAKEfakefakefakefakefakefakefakefakefake")))
	wrongPrivKeyByte, _ := secp256k1.FromECDSAPrivateKey(wrongPrivKey)
	wrongPubKey := append(wrongPrivKey.PublicKey.X.Bytes(), wrongPrivKey.PublicKey.Y.Bytes()...)
	//wrongPubKeyHash, _ := NewUserPubKeyHash(wrongPubKey)
	//wrongAddress := KeyPair{*wrongPrivKey, wrongPubKey}.GenerateAddress()
	utxoIndex := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage()))
	utxoTx := utxo.NewUTXOTx()

	utxoTx.PutUtxo(&utxo.UTXO{transaction_base.TXOutput{common.NewAmount(4), pubKeyHash, ""}, []byte{1}, 0, utxo.UtxoNormal})
	utxoTx.PutUtxo(&utxo.UTXO{transaction_base.TXOutput{common.NewAmount(3), pubKeyHash, ""}, []byte{2}, 1, utxo.UtxoNormal})

	utxoIndex.SetIndex(map[string]*utxo.UTXOTx{
		pubKeyHash.String(): &utxoTx,
	})

	// Prepare a transaction to be verified
	txin := []transaction_base.TXInput{{[]byte{1}, 0, nil, pubKey}}
	txin1 := append(txin, transaction_base.TXInput{[]byte{2}, 1, nil, pubKey})      // Normal test
	txin2 := append(txin, transaction_base.TXInput{[]byte{2}, 1, nil, wrongPubKey}) // previous not found with wrong pubkey
	txin3 := append(txin, transaction_base.TXInput{[]byte{3}, 1, nil, pubKey})      // previous not found with wrong Txid
	txin4 := append(txin, transaction_base.TXInput{[]byte{2}, 2, nil, pubKey})      // previous not found with wrong TxIndex
	pbh, _ := account.NewUserPubKeyHash(pubKey)
	txout := []transaction_base.TXOutput{{common.NewAmount(7), pbh, ""}}
	txout2 := []transaction_base.TXOutput{{common.NewAmount(8), pbh, ""}} //Vout amount > Vin amount

	tests := []struct {
		name     string
		tx       Transaction
		signWith []byte
		ok       error
	}{
		{"normal", Transaction{nil, txin1, txout, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}, privKeyByte, nil},
		{"previous tx not found with wrong pubkey", Transaction{nil, txin2, txout, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}, privKeyByte, errors.New("Transaction: prevUtxos not found")},
		{"previous tx not found with wrong Txid", Transaction{nil, txin3, txout, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}, privKeyByte, errors.New("Transaction: prevUtxos not found")},
		{"previous tx not found with wrong TxIndex", Transaction{nil, txin4, txout, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}, privKeyByte, errors.New("Transaction: prevUtxos not found")},
		{"Amount invalid", Transaction{nil, txin1, txout2, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}, privKeyByte, errors.New("Transaction: ID is invalid")},
		{"transaction_logic.Sign invalid", Transaction{nil, txin1, txout, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}, wrongPrivKeyByte, errors.New("Transaction: ID is invalid")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tx.ID = tt.tx.Hash()
			// Generate signatures for all tx inputs
			for i := range tt.tx.Vin {
				txCopy := tt.tx.TrimmedCopy(false)
				txCopy.Vin[i].Signature = nil
				txCopy.Vin[i].PubKey = []byte(pubKeyHash)
				signature, _ := secp256k1.Sign(txCopy.Hash(), tt.signWith)
				tt.tx.Vin[i].Signature = signature
			}

			// Verify the signatures
			_, err := transaction_logic.VerifyTransaction(utxoIndex, &tt.tx, 0)
			assert.Equal(t, tt.ok, err)
		})
	}
}

func TestInvalidExecutionTx(t *testing.T) {
	var prikey1 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa71"
	var pubkey1 = account.GenerateKeyPairByPrivateKey(prikey1).GetPublicKey()
	var pkHash1, _ = account.NewUserPubKeyHash(pubkey1)
	var deploymentTx = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{tx1.ID, 1, nil, pubkey1},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(5), pkHash1, "dapp_schedule"},
		},
		Tip: common.NewAmount(1),
	}
	deploymentTx.ID = deploymentTx.Hash()
	contractPubkeyHash := deploymentTx.Vout[0].PubKeyHash

	utxoIndex := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage()))
	utxoTx := utxo.NewUTXOTx()

	utxoTx.PutUtxo(&utxo.UTXO{deploymentTx.Vout[0], deploymentTx.ID, 0, utxo.UtxoNormal})
	utxoIndex.SetIndex(map[string]*utxo.UTXOTx{
		pkHash1.String(): &utxoTx,
	})

	var executionTx = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{deploymentTx.ID, 0, nil, pubkey1},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(3), contractPubkeyHash, "execution"},
		},
		Tip: common.NewAmount(2),
	}
	executionTx.ID = executionTx.Hash()
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey1).GetPrivateKey(), utxoIndex.GetAllUTXOsByPubKeyHash(pkHash1).GetAllUtxos(), &executionTx)

	_, err1 := transaction_logic.VerifyTransaction(utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage())), &executionTx, 0)
	_, err2 := transaction_logic.VerifyTransaction(utxoIndex, &executionTx, 0)
	assert.NotNil(t, err1)
	assert.Nil(t, err2)
}

func TestNewCoinbaseTX(t *testing.T) {
	t1 := transaction_logic.NewCoinbaseTX(account.NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 0, common.NewAmount(0))
	expectVin := transaction_base.TXInput{nil, -1, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte("Reward to 'dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB'")}
	expectVout := transaction_base.TXOutput{common.NewAmount(10000000), account.PubKeyHash([]byte{0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80, 0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49, 0xe5, 0x27, 0xf0, 0x42, 0x5d}), ""}
	assert.Equal(t, 1, len(t1.Vin))
	assert.Equal(t, expectVin, t1.Vin[0])
	assert.Equal(t, 1, len(t1.Vout))
	assert.Equal(t, expectVout, t1.Vout[0])
	assert.Equal(t, common.NewAmount(0), t1.Tip)

	t2 := transaction_logic.NewCoinbaseTX(account.NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 0, common.NewAmount(0))

	// Assert that transaction_logic.NewCoinbaseTX is deterministic (i.e. >1 coinbaseTXs in a block would have identical txid)
	assert.Equal(t, t1, t2)

	t3 := transaction_logic.NewCoinbaseTX(account.NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 1, common.NewAmount(0))

	assert.NotEqual(t, t1, t3)
	assert.NotEqual(t, t1.ID, t3.ID)
}

//test IsCoinBase function
func TestIsCoinBase(t *testing.T) {
	var tx1 = Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(2),
	}

	assert.False(t, tx1.IsCoinbase())

	t2 := transaction_logic.NewCoinbaseTX(account.NewAddress("13ZRUc4Ho3oK3Cw56PhE5rmaum9VBeAn5F"), "", 0, common.NewAmount(0))

	assert.True(t, t2.IsCoinbase())

}

func TestNewRewardTx(t *testing.T) {
	rewards := map[string]string{
		"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB": "8",
		"dastXXWLe5pxbRYFhcyUq8T3wb5srWkHKa": "9",
	}
	tx := NewRewardTx(5, rewards)

	values := []*common.Amount{tx.Vout[0].Value, tx.Vout[1].Value}
	assert.Contains(t, values, common.NewAmount(8))
	assert.Contains(t, values, common.NewAmount(9))
}

func TestTransaction_IsRewardTx(t *testing.T) {
	tests := []struct {
		name        string
		tx          Transaction
		expectedRes bool
	}{
		{"normal", NewRewardTx(1, map[string]string{"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB": "9"}), true},
		{"no rewards", NewRewardTx(1, nil), true},
		{"coinbase", transaction_logic.NewCoinbaseTX(account.NewAddress("dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB"), "", 5, common.NewAmount(0)), false},
		{"normal tx", *core.MockTransaction(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedRes, tt.tx.IsRewardTx())
		})
	}
}

func TestTransaction_Proto(t *testing.T) {
	tx1 := Transaction{
		ID:   util.GenerateRandomAoB(1),
		Vin:  GenerateFakeTxInputs(),
		Vout: GenerateFakeTxOutputs(),
		Tip:  common.NewAmount(5),
	}

	pb := tx1.ToProto()
	var i interface{} = pb
	_, correct := i.(proto.Message)
	assert.Equal(t, true, correct)
	mpb, err := proto.Marshal(pb)
	assert.Nil(t, err)

	newpb := &transactionpb.Transaction{}
	err = proto.Unmarshal(mpb, newpb)
	assert.Nil(t, err)

	tx2 := Transaction{}
	tx2.FromProto(newpb)

	assert.Equal(t, tx1, tx2)
}

func TestTransaction_GetContractAddress(t *testing.T) {

	tests := []struct {
		name        string
		addr        string
		expectedRes string
	}{
		{
			name:        "ContainsContractAddress",
			addr:        "cavQdWxvUQU1HhBg1d7zJFwhf31SUaQwop",
			expectedRes: "cavQdWxvUQU1HhBg1d7zJFwhf31SUaQwop",
		},
		{
			name:        "ContainsUserAddress",
			addr:        "dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expectedRes: "",
		},
		{
			name:        "EmptyInput",
			addr:        "",
			expectedRes: "",
		},
		{
			name:        "InvalidAddress",
			addr:        "dsdGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf",
			expectedRes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := account.NewAddress(tt.addr)
			pkh, _ := account.GeneratePubKeyHashByAddress(addr)
			tx := Transaction{
				nil,
				nil,
				[]transaction_base.TXOutput{
					{nil,
						pkh,
						"",
					},
				},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			}

			assert.Equal(t, account.NewAddress(tt.expectedRes), tx.GetContractAddress())
		})
	}
}

func TestTransaction_Execute(t *testing.T) {

	tests := []struct {
		name              string
		scAddr            string
		toAddr            string
		expectContractRun bool
	}{
		{
			name:              "CallAContract",
			scAddr:            "cWDSCWqwYRM6jNiN83PuRGvtcDuPpzBcfb",
			toAddr:            "cWDSCWqwYRM6jNiN83PuRGvtcDuPpzBcfb",
			expectContractRun: true,
		},
		{
			name:              "CallAWrongContractAddr",
			scAddr:            "cWDSCWqwYRM6jNiN83PuRGvtcDuPpzBcfb",
			toAddr:            "cavQdWxvUQU1HhBg1d7zJFwhf31SUaQwop",
			expectContractRun: false,
		},
		{
			name:              "NoPreviousContract",
			scAddr:            "",
			toAddr:            "cavQdWxvUQU1HhBg1d7zJFwhf31SUaQwop",
			expectContractRun: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := new(core.MockScEngine)
			contract := "helloworld!"
			toPKH, _ := account.GeneratePubKeyHashByAddress(account.NewAddress(tt.toAddr))
			scPKH, _ := account.GeneratePubKeyHashByAddress(account.NewAddress(tt.scAddr))

			scUtxo := utxo.UTXO{
				TxIndex: 0,
				Txid:    nil,
				TXOutput: transaction_base.TXOutput{
					PubKeyHash: scPKH,
					Contract:   contract,
				},
			}
			tx := ContractTx{Transaction{
				Vout:     []transaction_base.TXOutput{{nil, toPKH, "{\"function\":\"record\",\"args\":[\"dEhFf5mWTSe67mbemZdK3WiJh8FcCayJqm\",\"4\"]}"}},
				GasLimit: common.NewAmount(0),
				GasPrice: common.NewAmount(0),
			}}

			index := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage()))
			if tt.scAddr != "" {
				index.AddUTXO(scUtxo.TXOutput, nil, 0)
			}

			if tt.expectContractRun {
				sc.On("ImportSourceCode", contract)
				sc.On("ImportLocalStorage", mock.Anything)
				sc.On("ImportContractAddr", mock.Anything)
				sc.On("ImportUTXOs", mock.Anything)
				sc.On("ImportSourceTXID", mock.Anything)
				sc.On("ImportRewardStorage", mock.Anything)
				sc.On("ImportTransaction", mock.Anything)
				sc.On("ImportContractCreateUTXO", mock.Anything)
				sc.On("ImportPrevUtxos", mock.Anything)
				sc.On("GetGeneratedTXs").Return([]*Transaction{})
				sc.On("ImportCurrBlockHeight", mock.Anything)
				sc.On("ImportSeed", mock.Anything)
				sc.On("Execute", mock.Anything, mock.Anything).Return("")
			}
			parentBlk := core.GenerateMockBlock()
			preUTXO, err := utxo_logic.FindVinUtxosInUtxoPool(*index, tx.Transaction)

			if err != nil {
				println(err.Error())
			}
			isSCUTXO := (*index).GetAllUTXOsByPubKeyHash([]byte(tx.Vout[0].PubKeyHash)).Size() == 0
			transaction_logic.Execute(&tx, preUTXO, isSCUTXO, *index, core.NewScState(), nil, sc, 0, parentBlk)
			sc.AssertExpectations(t)
		})
	}
}

func TestTransaction_MatchRewards(t *testing.T) {

	tests := []struct {
		name          string
		tx            *Transaction
		rewardStorage map[string]string
		expectedRes   bool
	}{
		{"normal",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(1),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB": "1"},
			true,
		},
		{"emptyVout",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB": "1"},
			false,
		},
		{"emptyRewardMap",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(1),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			nil,
			false,
		},
		{"Wrong address",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(1),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{"dXnq2R6SzRNUt7ZsNAqyZc2P9ziF6vYekB": "1"},
			false,
		},
		{"Wrong amount",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(3),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{"dXnq2R6SzRNUt7ZsNAqyZc2P9ziF6vYekB": "1"},
			false,
		},
		{"twoAddresses",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(1),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				},
					{
						common.NewAmount(4),
						account.PubKeyHash([]byte{
							90, 13, 39, 130, 118, 11, 160, 130, 83, 126, 86, 102, 252, 178, 87,
							218, 57, 174, 123, 244, 229}),
						"",
					}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{
				"dEcqjSgREFi9gTCbAWpEQ3kbPxgsBzzhWS": "4",
				"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB": "1",
			},
			true,
		},
		{"MoreRewards",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(1),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				},
					{
						common.NewAmount(4),
						account.PubKeyHash([]byte{
							90, 13, 39, 130, 118, 11, 160, 130, 83, 126, 86, 102, 252, 178, 87,
							218, 57, 174, 123, 244, 229}),
						"",
					}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{
				"dEcqjSgREFi9gTCbAWpEQ3kbPxgsBzzhWS": "4",
				"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6vYekB": "1",
				"dXnq2R6SzRNUt7ZANAqyZc2P9ziF6fYekB": "3",
			},
			false,
		},
		{"MoreVout",
			&Transaction{
				nil,
				[]transaction_base.TXInput{{nil, -1, nil, RewardTxData}},
				[]transaction_base.TXOutput{{
					common.NewAmount(1),
					account.PubKeyHash([]byte{
						0x5a, 0xc9, 0x85, 0x37, 0x92, 0x37, 0x76, 0x80,
						0xb1, 0x31, 0xa1, 0xab, 0xb, 0x5b, 0xa6, 0x49,
						0xe5, 0x27, 0xf0, 0x42, 0x5d}),
					"",
				},
					{
						common.NewAmount(4),
						account.PubKeyHash([]byte{
							90, 13, 39, 130, 118, 11, 160, 130, 83, 126, 86, 102, 252, 178, 87,
							218, 57, 174, 123, 244, 229}),
						"",
					}},
				common.NewAmount(0),
				common.NewAmount(0),
				common.NewAmount(0),
			},
			map[string]string{
				"dEcqjSgREFi9gTCbAWpEQ3kbPxgsBzzhWS": "4",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedRes, tt.tx.MatchRewards(tt.rewardStorage))
		})
	}
}

func TestTransaction_VerifyDependentTransactions(t *testing.T) {
	var prikey1 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa71"
	var pubkey1 = account.GenerateKeyPairByPrivateKey(prikey1).GetPublicKey()
	var pkHash1, _ = account.NewUserPubKeyHash(pubkey1)
	var prikey2 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa72"
	var pubkey2 = account.GenerateKeyPairByPrivateKey(prikey2).GetPublicKey()
	var pkHash2, _ = account.NewUserPubKeyHash(pubkey2)
	var prikey3 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa73"
	var pubkey3 = account.GenerateKeyPairByPrivateKey(prikey3).GetPublicKey()
	var pkHash3, _ = account.NewUserPubKeyHash(pubkey3)
	var prikey4 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa74"
	var pubkey4 = account.GenerateKeyPairByPrivateKey(prikey4).GetPublicKey()
	var pkHash4, _ = account.NewUserPubKeyHash(pubkey4)
	var prikey5 = "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa75"
	var pubkey5 = account.GenerateKeyPairByPrivateKey(prikey5).GetPublicKey()
	var pkHash5, _ = account.NewUserPubKeyHash(pubkey5)

	var dependentTx1 = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{tx1.ID, 1, nil, pubkey1},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(5), pkHash1, ""},
			{common.NewAmount(10), pkHash2, ""},
		},
		Tip: common.NewAmount(3),
	}
	dependentTx1.ID = dependentTx1.Hash()

	var dependentTx2 = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{dependentTx1.ID, 1, nil, pubkey2},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(5), pkHash3, ""},
			{common.NewAmount(3), pkHash4, ""},
		},
		Tip: common.NewAmount(2),
	}
	dependentTx2.ID = dependentTx2.Hash()

	var dependentTx3 = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{dependentTx2.ID, 0, nil, pubkey3},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(1), pkHash4, ""},
		},
		Tip: common.NewAmount(4),
	}
	dependentTx3.ID = dependentTx3.Hash()

	var dependentTx4 = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{dependentTx2.ID, 1, nil, pubkey4},
			{dependentTx3.ID, 0, nil, pubkey4},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(3), pkHash1, ""},
		},
		Tip: common.NewAmount(1),
	}
	dependentTx4.ID = dependentTx4.Hash()

	var dependentTx5 = Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{dependentTx1.ID, 0, nil, pubkey1},
			{dependentTx4.ID, 0, nil, pubkey1},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(4), pkHash5, ""},
		},
		Tip: common.NewAmount(4),
	}
	dependentTx5.ID = dependentTx5.Hash()

	utxoIndex := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage()))

	utxoTx2 := utxo.NewUTXOTx()
	utxoTx2.PutUtxo(&utxo.UTXO{dependentTx1.Vout[1], dependentTx1.ID, 1, utxo.UtxoNormal})

	utxoTx1 := utxo.NewUTXOTx()
	utxoTx1.PutUtxo(&utxo.UTXO{dependentTx1.Vout[0], dependentTx1.ID, 0, utxo.UtxoNormal})

	utxoIndex.SetIndex(map[string]*utxo.UTXOTx{
		pkHash2.String(): &utxoTx2,
		pkHash1.String(): &utxoTx1,
	})

	tx2Utxo1 := utxo.UTXO{dependentTx2.Vout[0], dependentTx2.ID, 0, utxo.UtxoNormal}
	tx2Utxo2 := utxo.UTXO{dependentTx2.Vout[1], dependentTx2.ID, 1, utxo.UtxoNormal}
	tx2Utxo3 := utxo.UTXO{dependentTx3.Vout[0], dependentTx3.ID, 0, utxo.UtxoNormal}
	tx2Utxo4 := utxo.UTXO{dependentTx1.Vout[0], dependentTx1.ID, 0, utxo.UtxoNormal}
	tx2Utxo5 := utxo.UTXO{dependentTx4.Vout[0], dependentTx4.ID, 0, utxo.UtxoNormal}
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey2).GetPrivateKey(), utxoIndex.GetAllUTXOsByPubKeyHash(pkHash2).GetAllUtxos(), &dependentTx2)
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey3).GetPrivateKey(), []*utxo.UTXO{&tx2Utxo1}, &dependentTx3)
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey4).GetPrivateKey(), []*utxo.UTXO{&tx2Utxo2, &tx2Utxo3}, &dependentTx4)
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey1).GetPrivateKey(), []*utxo.UTXO{&tx2Utxo4, &tx2Utxo5}, &dependentTx5)

	txPool := core.NewTransactionPool(nil, 6000000)
	// verify dependent txs 2,3,4,5 with relation:
	//tx1 (UtxoIndex)
	//|     \
	//tx2    \
	//|  \    \
	//tx3-tx4-tx5

	// test a transaction whose Vin is from UtxoIndex
	_, err1 := transaction_logic.VerifyTransaction(utxoIndex, &dependentTx2, 0)
	assert.Nil(t, err1)
	txPool.Push(dependentTx2)

	// test a transaction whose Vin is from another transaction in transaction pool
	utxoIndex2 := *utxoIndex.DeepCopy()
	utxoIndex2.UpdateUtxoState(txPool.GetTransactions())
	_, err2 := transaction_logic.VerifyTransaction(&utxoIndex2, &dependentTx3, 0)
	assert.Nil(t, err2)
	txPool.Push(dependentTx3)

	// test a transaction whose Vin is from another two transactions in transaction pool
	utxoIndex3 := *utxoIndex.DeepCopy()
	utxoIndex3.UpdateUtxoState(txPool.GetTransactions())
	_, err3 := transaction_logic.VerifyTransaction(&utxoIndex3, &dependentTx4, 0)
	assert.Nil(t, err3)
	txPool.Push(dependentTx4)

	// test a transaction whose Vin is from another transaction in transaction pool and UtxoIndex
	utxoIndex4 := *utxoIndex.DeepCopy()
	utxoIndex4.UpdateUtxoState(txPool.GetTransactions())
	_, err4 := transaction_logic.VerifyTransaction(&utxoIndex4, &dependentTx5, 0)
	assert.Nil(t, err4)
	txPool.Push(dependentTx5)

	// test UTXOs not found for parent transactions
	_, err5 := transaction_logic.VerifyTransaction(utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage())), &dependentTx3, 0)
	assert.NotNil(t, err5)

	// test a standalone transaction
	txPool.Push(tx1)
	_, err6 := transaction_logic.VerifyTransaction(utxoIndex, &tx1, 0)
	assert.NotNil(t, err6)
}