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

package core

import (
	"github.com/dappley/go-dappley/network/network_model"
)

type Consensus interface {
	Validate(block *Block) bool

	Setup(NetService, string, *BlockChainManager)
	SetKey(string)

	// Start runs the consensus algorithm and begins to produce blocks
	Start()

	// Stop ceases the consensus algorithm and block production
	Stop()

	// IsProducingBlock returns true if this node itself is currently producing a block
	IsProducingBlock() bool

	// TODO: Should separate the concept of producers from PoW
	AddProducer(string) error
	GetProducers() []string
	//Return the lib block and new block whether pass lib policy
	CheckLibPolicy(b *Block) (*Block, bool)
}

type NetService interface {
	GetHostPeerInfo() *network_model.PeerInfo
}

type ScEngineManager interface {
	CreateEngine() ScEngine
	RunScheduledEvents(contractUtxo []*UTXO, scStorage *ScState, blkHeight uint64, seed int64)
}

type ScEngine interface {
	DestroyEngine()
	ImportSourceCode(source string)
	ImportLocalStorage(state *ScState)
	ImportContractAddr(contractAddr Address)
	ImportSourceTXID(txid []byte)
	ImportUTXOs(utxos []*UTXO)
	ImportRewardStorage(rewards map[string]string)
	ImportTransaction(tx *Transaction)
	ImportContractCreateUTXO(utxo *UTXO)
	ImportPrevUtxos(utxos []*UTXO)
	ImportCurrBlockHeight(currBlkHeight uint64)
	ImportSeed(seed int64)
	ImportNodeAddress(addr Address)
	GetGeneratedTXs() []*Transaction
	Execute(function, args string) string
}
