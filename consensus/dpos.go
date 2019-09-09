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
	"github.com/dappley/go-dappley/core/block_producer_info"
	"strings"
	"time"

	"github.com/dappley/go-dappley/core/block"
	"github.com/dappley/go-dappley/logic/block_logic"

	"github.com/dappley/go-dappley/core/account"
	"github.com/dappley/go-dappley/crypto/keystore/secp256k1"
	"github.com/hashicorp/golang-lru"
	logger "github.com/sirupsen/logrus"
)

const (
	MinConsensusSize = 4
)

type DPOS struct {
	producer        *block_producer_info.BlockProducerInfo
	producerKey     string
	stopCh          chan bool
	dynasty         *Dynasty
	slot            *lru.Cache
	lastProduceTime int64
}

//NewDPOS returns a new DPOS instance
func NewDPOS(producer *block_producer_info.BlockProducerInfo) *DPOS {
	dpos := &DPOS{
		producer: producer,
		stopCh:   make(chan bool, 1),
	}

	slot, err := lru.New(128)
	if err != nil {
		logger.Panic(err)
	}
	dpos.slot = slot
	return dpos
}

//SetKey sets the producer key
func (dpos *DPOS) SetKey(key string) {
	dpos.producerKey = key
}

//SetDynasty sets the dynasty
func (dpos *DPOS) SetDynasty(dynasty *Dynasty) {
	dpos.dynasty = dynasty
}

//GetDynasty returns the dynasty
func (dpos *DPOS) GetDynasty() *Dynasty {
	return dpos.dynasty
}

//AddProducer adds a producer to the dynasty
func (dpos *DPOS) AddProducer(producer string) error {
	err := dpos.dynasty.AddProducer(producer)
	return err
}

//GetProducers returns all current producers
func (dpos *DPOS) GetProducers() []string {
	return dpos.dynasty.GetProducers()
}

//GetProducerAddress returns the local producer's address
func (dpos *DPOS) GetProducerAddress() string {
	return dpos.producer.Beneficiary()
}

//Stop stops the current produce block process
func (dpos *DPOS) Stop() {
	logger.Info("DPoS stops...")
	dpos.stopCh <- true
}

//ProduceBlock starts producing block according to dpos consensus
func (dpos *DPOS) ProduceBlock(ProduceBlockFunc func(process func(*block.Block))) {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case now := <-ticker:
			if dpos.dynasty.IsMyTurn(dpos.producer.Beneficiary(), now.Unix()) {
				ProduceBlockFunc(dpos.hashAndSign)
				return
			}
		case <-dpos.stopCh:
			return
		}
	}
}

//Produced returns if the local producer produced the block
func (dpos *DPOS) Produced(blk *block.Block) bool {
	if blk != nil {
		return dpos.producer.Produced(blk)
	}
	return false
}

//hashAndSign signs the block
func (dpos *DPOS) hashAndSign(blk *block.Block) {
	hash := block_logic.CalculateHash(blk)
	blk.SetHash(hash)
	ok := block_logic.SignBlock(blk, dpos.producerKey)
	if !ok {
		logger.Warn("DPoS: failed to sign the new block.")
	}
}

// Validate checks that the block fulfills the dpos requirement and accepts the block in the time slot
func (dpos *DPOS) Validate(block *block.Block) bool {
	producerIsValid := dpos.verifyProducer(block)
	if !producerIsValid {
		return false
	}

	if dpos.isDoubleMint(block) {
		logger.Warn("DPoS: double-minting is detected.")
		return false
	}

	dpos.cacheBlock(block)
	return true
}

// verifyProducer verifies a given block is produced by the valid producer by verifying the signature of the block
func (dpos *DPOS) verifyProducer(block *block.Block) bool {
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

	pubKeyHash, err := account.NewUserPubKeyHash(pubkey[1:])
	if err != nil {
		logger.WithError(err).Warn("DPoS: cannot compute the public key hash!")
		return false
	}

	address := pubKeyHash.GenerateAddress()

	if strings.Compare(address.String(), producer) != 0 {
		logger.Warn("DPoS: the signer is not the producer in this time slot.")
		return false
	}

	if !dpos.isProducerBeneficiary(block) {
		logger.Debug("DPoS: failed to validate producer.")
		return false
	}

	return true
}

// isProducerBeneficiary is a requirement that ensures the reward is paid to the producer at the time slot
func (dpos *DPOS) isProducerBeneficiary(block *block.Block) bool {
	if block == nil {
		logger.Debug("DPoS: block is empty.")
		return false
	}

	producer := dpos.dynasty.ProducerAtATime(block.GetTimestamp())
	producerHash, _ := account.GeneratePubKeyHashByAddress(account.NewAddress(producer))

	cbtx := block.GetCoinbaseTransaction()
	if cbtx == nil {
		logger.Debug("DPoS: coinbase tx is empty.")
		return false
	}

	if len(cbtx.Vout) == 0 {
		logger.Debug("DPoS: coinbase vout is empty.")
		return false
	}

	return bytes.Compare(producerHash, cbtx.Vout[0].PubKeyHash) == 0
}

//isDoubleMint returns if the block's producer has already produced a block in the current time slot
func (dpos *DPOS) isDoubleMint(blk *block.Block) bool {
	existBlock, exist := dpos.slot.Get(int(blk.GetTimestamp() / int64(dpos.GetDynasty().timeBetweenBlk)))
	if !exist {
		return false
	}

	return !block_logic.IsHashEqual(existBlock.(*block.Block).GetHash(), blk.GetHash())
}

//cacheBlock adds the block to cache for double minting check
func (dpos *DPOS) cacheBlock(block *block.Block) {
	dpos.slot.Add(int(block.GetTimestamp()/int64(dpos.GetDynasty().timeBetweenBlk)), block)
}

//GetLibProducerNum returns the minimum number of producers required
func (dpos *DPOS) GetLibProducerNum() int {
	return len(dpos.dynasty.GetProducers())*2/3 + 1
}

//IsBypassingLibCheck returns if LIB check should be skipped
func (dpos *DPOS) IsBypassingLibCheck() bool {
	return len(dpos.dynasty.GetProducers()) < MinConsensusSize
}

//IsNonRepeatingBlockProducerRequired returns if the same producer can producer consecutive blocks
func (dpos *DPOS) IsNonRepeatingBlockProducerRequired() bool {
	return true
}
