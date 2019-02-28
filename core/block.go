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
// You should have received a copy of the GNU Gc
// eneral Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//
package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/core/pb"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
	"github.com/dappley/go-dappley/crypto/sha3"
	"github.com/dappley/go-dappley/util"
)

type BlockHeader struct {
	hash      Hash
	prevHash  Hash
	nonce     int64
	timestamp int64
	sign      Hash
	height    uint64
}

type Block struct {
	header       *BlockHeader
	transactions []*Transaction
}

type Hash []byte

func (h Hash) String() string {
	return hex.EncodeToString(h)
}

func NewBlock(txs []*Transaction, parent *Block) *Block {
	return NewBlockWithTimestamp(txs, parent, time.Now().Unix())
}

func NewBlockWithTimestamp(txs []*Transaction, parent *Block, timeStamp int64) *Block {

	var prevHash []byte
	var height uint64
	height = 1
	if parent != nil {
		prevHash = parent.GetHash()
		height = parent.GetHeight() + 1
	}

	if txs == nil {
		txs = []*Transaction{}
	}
	return &Block{
		header: &BlockHeader{
			hash:      []byte{},
			prevHash:  prevHash,
			nonce:     0,
			timestamp: timeStamp,
			sign:      nil,
			height:    height,
		},
		transactions: txs,
	}
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.transactions {
		txHashes = append(txHashes, tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) Serialize() []byte {
	rawBytes, err := proto.Marshal(b.ToProto())
	if err!= nil {
		logger.WithError(err).Panic("Block: Cannot serialize block!")
	}
	return rawBytes
}

func Deserialize(d []byte) *Block {
	pb := &corepb.Block{}
	err := proto.Unmarshal(d, pb)
	if err!= nil {
		logger.WithError(err).Panic("Block: Cannot deserialize block!")
	}
	block := &Block{}
	block.FromProto(pb)
	return block
}

func (b *Block) GetHeader() *BlockHeader {
	return b.header
}

func (b *Block) SetHash(hash Hash) {
	b.header.hash = hash
}

func (b *Block) GetHash() Hash {
	return b.header.hash
}

func (b *Block) GetSign() Hash {
	return b.header.sign
}

func (b *Block) GetHeight() uint64 {
	return b.header.height
}

func (b *Block) GetPrevHash() Hash {
	return b.header.prevHash
}

func (b *Block) SetNonce(nonce int64) {
	b.header.nonce = nonce
}

func (b *Block) GetNonce() int64 {
	return b.header.nonce
}

func (b *Block) GetTimestamp() int64 {
	return b.header.timestamp
}

func (b *Block) GetTransactions() []*Transaction {
	return b.transactions
}

func (b *Block) ToProto() proto.Message {

	var txArray []*corepb.Transaction
	for _, tx := range b.transactions {
		txArray = append(txArray, tx.ToProto().(*corepb.Transaction))
	}

	return &corepb.Block{
		Header:       b.header.ToProto().(*corepb.BlockHeader),
		Transactions: txArray,
	}
}

func (b *Block) FromProto(pb proto.Message) {

	bh := BlockHeader{}
	bh.FromProto(pb.(*corepb.Block).GetHeader())
	b.header = &bh

	var txs []*Transaction

	for _, txpb := range pb.(*corepb.Block).GetTransactions() {
		tx := &Transaction{}
		tx.FromProto(txpb)
		txs = append(txs, tx)
	}
	b.transactions = txs
}

func (bh *BlockHeader) ToProto() proto.Message {
	return &corepb.BlockHeader{
		Hash:         bh.hash,
		PreviousHash: bh.prevHash,
		Nonce:        bh.nonce,
		Timestamp:    bh.timestamp,
		Signature:    bh.sign,
		Height:       bh.height,
	}
}

func (bh *BlockHeader) FromProto(pb proto.Message) {
	if pb == nil {
		return
	}
	bh.hash = pb.(*corepb.BlockHeader).GetHash()
	bh.prevHash = pb.(*corepb.BlockHeader).GetPreviousHash()
	bh.nonce = pb.(*corepb.BlockHeader).GetNonce()
	bh.timestamp = pb.(*corepb.BlockHeader).GetTimestamp()
	bh.sign = pb.(*corepb.BlockHeader).GetSignature()
	bh.height = pb.(*corepb.BlockHeader).GetHeight()
}

func (b *Block) CalculateHash() Hash {
	return b.CalculateHashWithNonce(b.GetNonce())
}

func (b *Block) CalculateHashWithoutNonce() Hash {
	data := bytes.Join(
		[][]byte{
			b.GetPrevHash(),
			b.HashTransactions(),
			util.IntToHex(b.GetTimestamp()),
		},
		[]byte{},
	)

	hasher := sha3.New256()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func (b *Block) CalculateHashWithNonce(nonce int64) Hash {
	data := bytes.Join(
		[][]byte{
			b.GetPrevHash(),
			b.HashTransactions(),
			util.IntToHex(b.GetTimestamp()),
			//util.IntToHex(targetBits),
			util.IntToHex(nonce),
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) SignBlock(key string, data []byte) bool {
	if len(key) <= 0 {
		logger.Warn("Block: the key is too short for signature!")
		return false
	}
	privData, err := hex.DecodeString(key)

	if err != nil {
		logger.Warn("Block: cannot decode private key for signature!")
		return false
	}
	signature, err := secp256k1.Sign(data, privData)
	if err != nil {
		logger.WithError(err).Warn("Block: failed to calculate signature!")
		return false
	}

	b.header.sign = signature
	return true
}

func (b *Block) VerifyHash() bool {
	return bytes.Compare(b.GetHash(), b.CalculateHash()) == 0
}

func (b *Block) VerifyTransactions(utxo UTXOIndex, scState *ScState, parentBlk *Block, txPool *TransactionPool) bool {
	if len(b.GetTransactions()) == 0 {
		logger.WithFields(logger.Fields{
			"hash": b.GetHash(),
		}).Debug("Block: there is no transaction to verify in this block.")
		return true
	}

	var rewardTX *Transaction

	for _, tx := range b.GetTransactions() {
		// Collect the contract-incurred transactions in this block
		if tx.IsRewardTx() {
			if rewardTX != nil {
				logger.Warn("Block: contains more than 1 reward transaction.")
				return false
			}
			rewardTX = tx
		} else if tx.IsFromContract() {
			txInPool := txPool.GetTransactionById(tx.ID)
			if !tx.IsIdentical(utxo, txInPool) {
				logger.Warn("Block: generated tx cannot be verified.")
				return false
			}
		} else {
			// tx is a normal transactions
			if !tx.Verify(&utxo, b.GetHeight()) {
				return false
			}
		}
		utxo.UpdateUtxo(tx)
	}

	return true
}

func IsParentBlockHash(parentBlk, childBlk *Block) bool {
	if parentBlk == nil || childBlk == nil {
		return false
	}
	return reflect.DeepEqual(parentBlk.GetHash(), childBlk.GetPrevHash())
}

func IsHashEqual(h1 Hash, h2 Hash) bool {

	return reflect.DeepEqual(h1, h2)
}

func IsParentBlockHeight(parentBlk, childBlk *Block) bool {
	if parentBlk == nil || childBlk == nil {
		return false
	}
	return parentBlk.GetHeight() == childBlk.GetHeight()-1
}

func (b *Block) IsParentBlock(child *Block) bool {
	return IsParentBlockHash(b, child) && IsParentBlockHeight(b, child)
}

func (b *Block) Rollback(txPool *TransactionPool) {
	if b != nil {
		for _, tx := range b.GetTransactions() {
			if !tx.IsCoinbase() && !tx.IsRewardTx() {
				txPool.Push(tx)
			}
		}
	}
}

func (b *Block) FindTransactionById(txid []byte) *Transaction {
	for _, tx := range b.transactions {
		if bytes.Compare(tx.ID, txid) == 0 {

			return tx
		}

	}
	return nil
}

func (b *Block) GetCoinbaseTransaction() *Transaction {
	//the coinbase transaction is usually placed at the end of all transactions
	for i := len(b.transactions) - 1; i >= 0; i-- {
		if b.transactions[i].IsCoinbase() {
			return b.transactions[i]
		}
	}
	return nil
}
