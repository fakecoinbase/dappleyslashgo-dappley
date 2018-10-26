// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package consensus

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Dpos struct {
	bc          *core.Blockchain
	delegate    BlockProducer
	mintBlkCh   chan *MinedBlock
	blkProduced bool
	node        core.NetService
	quitCh      chan bool
	dynasty     *Dynasty
	slot        *lru.Cache
}

func NewDpos() *Dpos {
	dpos := &Dpos{
		delegate:    NewDelegate(),
		mintBlkCh:   make(chan *MinedBlock, 1),
		blkProduced: false,
		node:        nil,
		quitCh:      make(chan bool, 1),
	}

	slot, err := lru.New(128)
	if err != nil {
		logger.Panic(err)
	}
	dpos.slot = slot
	return dpos
}

func (dpos *Dpos) GetSlot() *lru.Cache {
	return dpos.slot
}

func (dpos *Dpos) Setup(node core.NetService, cbAddr string) {
	dpos.bc = node.GetBlockchain()
	dpos.node = node
	dpos.delegate.Setup(dpos.bc, cbAddr, dpos.mintBlkCh)
}

func (dpos *Dpos) SetTargetBit(bit int) {
	//dpos.delegate.SetTargetBit(bit)
}

func (dpos *Dpos) SetKey(key string) {
	dpos.delegate.SetPrivKey(key)
}

func (dpos *Dpos) SetDynasty(dynasty *Dynasty) {
	dpos.dynasty = dynasty
}

func (dpos *Dpos) GetDynasty() *Dynasty {
	return dpos.dynasty
}

func (dpos *Dpos) AddProducer(producer string) error {
	err := dpos.dynasty.AddProducer(producer)
	return err
}

func (dpos *Dpos) GetProducers() []string {
	return dpos.dynasty.GetProducers()
}

func (dpos *Dpos) GetBlockChain() *core.Blockchain {
	return dpos.bc
}

func (dpos *Dpos) Validate(block *core.Block) bool {
	if !dpos.delegate.Validate(block) {
		logger.Debug("Dpos: blockProducer validate block failed")
		return false
	}
	if !dpos.dynasty.ValidateProducer(block) {
		logger.Debug("Dpos: producer validate failed")
		return false
	}
	if dpos.isDoubleMint(block) {
		logger.Debug("Dpos: doubleminting case found!")
		return false
	}

	dpos.slot.Add(block.GetTimestamp(), block)

	return true
}

func (dpos *Dpos) Start() {
	go func() {
		logger.Info("Dpos Starts...", dpos.node.GetPeerID())
		ticker := time.NewTicker(time.Second).C
		for {
			select {
			case now := <-ticker:
				if dpos.dynasty.IsMyTurn(dpos.delegate.Beneficiary(), now.Unix()) {
					logger.WithFields(logger.Fields{
						"peerid": dpos.node.GetPeerID(),
					}).Info("My Turn to Mint")
					dpos.blkProduced = false
					dpos.delegate.Start()
				}
			case minedBlk := <-dpos.mintBlkCh:
				dpos.blkProduced = true
				if minedBlk.isValid {
					logger.WithFields(logger.Fields{
						"peerid": dpos.node.GetPeerID(),
						"hash":   hex.EncodeToString(minedBlk.block.GetHash()),
					}).Info("Dpos: A Block has been mined!")
					dpos.updateNewBlock(minedBlk.block)
				}
			case <-dpos.quitCh:
				logger.WithFields(logger.Fields{
					"peerid": dpos.node.GetPeerID(),
				}).Info("Dpos: Dpos Stops!")
				return
			}
		}
	}()
}

func (dpos *Dpos) Stop() {
	dpos.quitCh <- true
	dpos.delegate.Stop()
}

func (dpos *Dpos) isForking() bool {
	return false
}

func (dpos *Dpos) isDoubleMint(block *core.Block) bool {
	if _, exist := dpos.slot.Get(block.GetTimestamp()); exist {
		logger.Debug("Someone is minting when they are not supposed to!")
		return true
	}
	return false
}
func (dpos *Dpos) StartNewBlockMinting() {
	dpos.delegate.Stop()
}
func (dpos *Dpos) FinishedMining() bool {
	return dpos.blkProduced
}

func (dpos *Dpos) updateNewBlock(newBlock *core.Block) {
	logger.WithFields(logger.Fields{
		"height": newBlock.GetHeight(),
		"hash":   hex.EncodeToString(newBlock.GetHash()),
	}).Info("DpoS: Minted a new block")
	dpos.bc.AddBlockToTail(newBlock)
	dpos.node.BroadcastBlock(newBlock)
}

func (dpos *Dpos) VerifyBlock(block *core.Block) bool {
	hash1 := block.GetHash()
	sign := block.GetSign()

	producer := dpos.dynasty.ProducerAtATime(block.GetTimestamp())

	if hash1 == nil {
		logger.Warn("DPoS: block hash empty!")
		return false
	}
	if sign == nil {
		logger.Warn("DPoS: block signature empty!")
		return false
	}

	pubkey, err := secp256k1.RecoverECDSAPublicKey(hash1, sign)
	if err != nil {
		logger.Warn("DPoS: Get pub key from block signature error!")
		return false
	}

	pubKeyHash,err := core.NewUserPubKeyHash(pubkey[1:])
	if err != nil {
		logger.Warn("DPoS: Invalid Public Key!")
		return false
	}

	address := pubKeyHash.GenerateAddress()

	if strings.Compare(address.String(), producer) == 0 {
		return true
	}

	logger.Warn("DPoS: Address is not current producer's")
	return false

}
