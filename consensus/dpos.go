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
	"bytes"
	"encoding/hex"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
)

const (
	ConsensusSize = defaultMaxProducers*2/3 + 1
)

type DPOS struct {
	bp          *BlockProducer
	producerKey string
	newBlockCh  chan *core.Block
	node        core.NetService
	stopCh      chan bool
	stopLibCh   chan bool
	dynasty     *Dynasty
	slot        *lru.Cache
}

func NewDPOS() *DPOS {
	dpos := &DPOS{
		bp:         NewBlockProducer(),
		newBlockCh: make(chan *core.Block, 1),
		node:       nil,
		stopCh:     make(chan bool, 1),
		stopLibCh:  make(chan bool, 1),
	}

	slot, err := lru.New(128)
	if err != nil {
		logger.Panic(err)
	}
	dpos.slot = slot
	return dpos
}

func (dpos *DPOS) GetSlot() *lru.Cache {
	return dpos.slot
}

func (dpos *DPOS) AddBlockToSlot(block *core.Block) {
	dpos.slot.Add(int(block.GetTimestamp()/int64(dpos.GetDynasty().timeBetweenBlk)), block)
}

func (dpos *DPOS) Setup(node core.NetService, cbAddr string) {
	dpos.node = node
	dpos.bp.Setup(node.GetBlockchain(), cbAddr)
	dpos.bp.SetProcess(dpos.hashAndSign)
}

func (dpos *DPOS) SetKey(key string) {
	dpos.producerKey = key
}

func (dpos *DPOS) SetDynasty(dynasty *Dynasty) {
	dpos.dynasty = dynasty
}

func (dpos *DPOS) GetDynasty() *Dynasty {
	return dpos.dynasty
}

func (dpos *DPOS) AddProducer(producer string) error {
	err := dpos.dynasty.AddProducer(producer)
	return err
}

func (dpos *DPOS) GetProducers() []string {
	return dpos.dynasty.GetProducers()
}

// Validate checks that the block fulfills the dpos requirement and accepts the block in the time slot
func (dpos *DPOS) Validate(block *core.Block) bool {
	producerIsValid := dpos.verifyProducer(block)
	if !producerIsValid {
		return false
	}

	if !dpos.beneficiaryIsProducer(block) {
		logger.Debug("DPoS: failed to validate producer.")
		return false
	}
	if dpos.isDoubleMint(block) {
		logger.Warn("DPoS: double-minting is detected.")
		return false
	}

	dpos.AddBlockToSlot(block)
	return true
}

func (dpos *DPOS) Start() {
	go func() {
		logger.WithFields(logger.Fields{
			"peer_id": dpos.node.GetPeerID(),
		}).Info("DPoS starts...")
		if len(dpos.stopCh) > 0 {
			<-dpos.stopCh
		}
		ticker := time.NewTicker(time.Second).C
		for {
			select {
			case now := <-ticker:
				if dpos.dynasty.IsMyTurn(dpos.bp.Beneficiary(), now.Unix()) {
					logger.WithFields(logger.Fields{
						"peer_id": dpos.node.GetPeerID(),
					}).Info("DPoS: it is my turn to produce block.")
					// Do not produce block if block pool is syncing
					if dpos.node.GetBlockchain().GetState() != core.BlockchainReady {
						logger.Info("DPoS: block producer paused because block pool is syncing.")
						continue
					}
					ctx := dpos.bp.ProduceBlock()
					if ctx == nil || !dpos.Validate(ctx.Block) {
						logger.Error("DPoS: produced an invalid block!")
						continue
					}
					dpos.updateNewBlock(ctx)
				}
			case <-dpos.stopCh:
				dpos.stopLibCh <- true
				return

			}
		}
	}()

	go func() {
		if len(dpos.stopLibCh) > 0 {
			<-dpos.stopLibCh
		}
		updateLibTicker := time.NewTicker(time.Second * defaultTimeBetweenBlk).C
		for {
			select {
			case <-updateLibTicker:
				dpos.UpdateLIB()
			case <-dpos.stopLibCh:
				return
			}
		}
	}()
}

func (dpos *DPOS) Stop() {
	logger.WithFields(logger.Fields{
		"peer_id": dpos.node.GetPeerID(),
	}).Info("DPoS stops...")
	dpos.stopCh <- true
}

func (dpos *DPOS) hashAndSign(ctx *core.BlockContext) {
	//block.SetNonce(0)
	hash := ctx.Block.CalculateHash()
	ctx.Block.SetHash(hash)
	ok := ctx.Block.SignBlock(dpos.producerKey, hash)
	if !ok {
		logger.Warn("DPoS: failed to sign the new block.")
	}
}

func (dpos *DPOS) isForking() bool {
	return false
}

func (dpos *DPOS) isDoubleMint(block *core.Block) bool {
	existBlock, exist := dpos.slot.Get(int(block.GetTimestamp() / int64(dpos.GetDynasty().timeBetweenBlk)))
	if !exist {
		return false
	}

	return !core.IsHashEqual(existBlock.(*core.Block).GetHash(), block.GetHash())
}

// verifyProducer verifies a given block is produced by the valid producer by verifying the signature of the block
func (dpos *DPOS) verifyProducer(block *core.Block) bool {
	if block == nil {
		logger.Warn("DPoS: block is empty!")
		return false
	}

	hash := block.GetHash()
	sign := block.GetSign()

	producer := dpos.dynasty.ProducerAtATime(block.GetTimestamp())

	if hash == nil {
		logger.Warn("DPoS: block hash is empty!")
		return false
	}
	if sign == nil {
		logger.Warn("DPoS: block signature is empty!")
		return false
	}

	pubkey, err := secp256k1.RecoverECDSAPublicKey(hash, sign)
	if err != nil {
		logger.WithError(err).Warn("DPoS: cannot recover the public key from the block signature!")
		return false
	}

	pubKeyHash, err := core.NewUserPubKeyHash(pubkey[1:])
	if err != nil {
		logger.WithError(err).Warn("DPoS: cannot compute the public key hash!")
		return false
	}

	address := pubKeyHash.GenerateAddress()

	if strings.Compare(address.String(), producer) != 0 {
		logger.Warn("DPoS: the signer is not the producer in this time slot.")
		return false
	}

	return true
}

// beneficiaryIsProducer is a requirement that ensures the reward is paid to the producer at the time slot
func (dpos *DPOS) beneficiaryIsProducer(block *core.Block) bool {
	if block == nil {
		logger.Debug("DPoS: block is empty.")
		return false
	}

	producer := dpos.dynasty.ProducerAtATime(block.GetTimestamp())
	producerHash, _ := core.NewAddress(producer).GetPubKeyHash()

	cbtx := block.GetCoinbaseTransaction()
	if cbtx == nil {
		logger.Debug("DPoS: coinbase tx is empty.")
		return false
	}

	if len(cbtx.Vout) == 0 {
		logger.Debug("DPoS: coinbase vout is empty.")
		return false
	}

	return bytes.Compare(producerHash, []byte(cbtx.Vout[0].PubKeyHash)) == 0
}

func (dpos *DPOS) IsProducingBlock() bool {
	return !dpos.bp.IsIdle()
}

func (dpos *DPOS) updateNewBlock(ctx *core.BlockContext) {
	logger.WithFields(logger.Fields{
		"peer_id": dpos.node.GetPeerID(),
		"height":  ctx.Block.GetHeight(),
		"hash":    hex.EncodeToString(ctx.Block.GetHash()),
	}).Info("DPoS: produced a new block.")
	if !ctx.Block.VerifyHash() {
		logger.Warn("DPoS: hash of the new block is invalid.")
		return
	}
	err := dpos.node.GetBlockchain().AddBlockContextToTail(ctx)
	if err != nil {
		logger.Warn(err)
		return
	}
	dpos.node.BroadcastBlock(ctx.Block)
}

func (dpos *DPOS) UpdateLIB() {
	bc := dpos.node.GetBlockchain()
	lib, _ := bc.GetLIB()
	tailBlock, _ := bc.GetTailBlock()
	cur := tailBlock
	miners := make(map[string]bool)

	for !cur.GetHash().Equals(lib.GetHash()) {
		producer := dpos.dynasty.ProducerAtATime(cur.GetTimestamp())
		if int(cur.GetHeight())-int(lib.GetHeight()) < ConsensusSize-len(miners) {
			logger.WithFields(logger.Fields{
				"miners.limit":     ConsensusSize,
				"miners.supported": len(miners),
			}).Info("DPoS: failed to update latest irreversible block.")
			return
		}
		miners[producer] = true
		if len(miners) >= ConsensusSize {
			if err := dpos.node.GetBlockchain().SetLIBHash(cur.GetHash()); err != nil {
				logger.WithFields(logger.Fields{
					"tail_hash": tailBlock.GetHash(),
					"lib_hash":  cur.GetHash(),
				}).Info("DPoS: failed to store latest irreversible block.")
				return
			}
			logger.WithFields(logger.Fields{
				"new lib_hash":     cur.GetHash(),
				"old lib_hash":     lib.GetHash(),
				"tail hash":        tailBlock.GetHash(),
				"miners limit":     ConsensusSize,
				"miners supported": len(miners),
			}).Info("DPoS: Succeed to update latest irreversible block.")
			return
		}
		tmp := cur
		cur, _ = dpos.node.GetBlockchain().GetBlockByHash(cur.GetPrevHash())
		if cur == nil || core.CheckGenesisBlock(cur) {
			logger.WithFields(logger.Fields{
				"tail": tailBlock,
				"cur":  tmp,
			}).Info("DPoS: failed to find latest irreversible block.")
			return
		}
	}
	logger.WithFields(logger.Fields{
		"cur":              cur.GetHash(),
		"lib":              lib.GetHash(),
		"tail":             tailBlock.GetHash(),
		"err":              "supported miners is not enough",
		"miners.limit":     ConsensusSize,
		"miners.supported": len(miners),
	}).Info("DPoS: failed to update latest irreversible block.")
}
