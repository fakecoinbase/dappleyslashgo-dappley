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
package intergationTest

import (
	"github.com/dappley/go-dappley/logic/block_producer"
	"reflect"
	"testing"
	"time"

	"github.com/dappley/go-dappley/core/block"
	"github.com/dappley/go-dappley/core/block_producer_info"
	"github.com/dappley/go-dappley/core/blockchain"
	"github.com/dappley/go-dappley/core/transaction"
	"github.com/dappley/go-dappley/core/transaction_base"
	"github.com/dappley/go-dappley/core/utxo"
	"github.com/dappley/go-dappley/logic"
	"github.com/dappley/go-dappley/logic/account_logic"
	"github.com/dappley/go-dappley/logic/block_logic"
	"github.com/dappley/go-dappley/logic/blockchain_logic"
	"github.com/dappley/go-dappley/logic/transaction_logic"
	"github.com/dappley/go-dappley/logic/transaction_pool"
	"github.com/dappley/go-dappley/logic/utxo_logic"

	"github.com/dappley/go-dappley/util"

	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/consensus"
	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/core/account"
	"github.com/dappley/go-dappley/network"
	"github.com/dappley/go-dappley/storage"
	"github.com/stretchr/testify/assert"
)

const testport_msg_relay_port1 = 31202
const testport_msg_relay_port2 = 31212
const testport_msg_relay_port3 = 31222
const testport_fork = 10500
const testport_fork_segment = 10511
const testport_fork_syncing = 10531
const testport_fork_download = 10600
const InvalidAddress = "Invalid Address"

//test logic.Send
func TestSend(t *testing.T) {
	var mineReward = common.NewAmount(10000000)
	testCases := []struct {
		name             string
		transferAmount   *common.Amount
		tipAmount        *common.Amount
		contract         string
		gasLimit         *common.Amount
		gasPrice         *common.Amount
		expectedTransfer *common.Amount
		expectedTip      *common.Amount
		expectedErr      error
	}{
		{"Deploy contract", common.NewAmount(7), common.NewAmount(0), "dapp_schedule!", common.NewAmount(30000), common.NewAmount(1), common.NewAmount(7), common.NewAmount(0), nil},
		{"Send with no tip", common.NewAmount(7), common.NewAmount(0), "", common.NewAmount(0), common.NewAmount(0), common.NewAmount(7), common.NewAmount(0), nil},
		{"Send with tips", common.NewAmount(6), common.NewAmount(2), "", common.NewAmount(0), common.NewAmount(0), common.NewAmount(6), common.NewAmount(2), nil},
		{"Send zero with no tip", common.NewAmount(0), common.NewAmount(0), "", common.NewAmount(0), common.NewAmount(0), common.NewAmount(0), common.NewAmount(0), logic.ErrInvalidAmount},
		{"Send zero with tips", common.NewAmount(0), common.NewAmount(2), "", common.NewAmount(0), common.NewAmount(0), common.NewAmount(0), common.NewAmount(0), logic.ErrInvalidAmount},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			store := storage.NewRamStorage()
			defer store.Close()

			// Create a account address
			SenderAccount, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
			if err != nil {
				panic(err)
			}
			// Create a miner account; Balance is 0 initially
			minerAccount, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
			if err != nil {
				panic(err)
			}

			node := network.FakeNodeWithPidAndAddr(store, "test", "test")
			// Create a PoW blockchain with the logic.Sender wallet's address as the coinbase address
			// i.e. logic.Sender's wallet would have mineReward amount after blockchain created
			bc, pow := CreateBlockchain(minerAccount.GetKeyPair().GenerateAddress(), SenderAccount.GetKeyPair().GenerateAddress(), store, transaction_pool.NewTransactionPool(node, 128))

			// Create a receiver account; Balance is 0 initially
			receiverAccount, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
			if err != nil {
				panic(err)
			}

			// logic.Send coins from logic.SenderAccount to receiverAccount
			var rcvAddr account.Address
			isContract := (tc.contract != "")
			if isContract {
				rcvAddr = account.NewAddress("")
			} else {
				rcvAddr = receiverAccount.GetKeyPair().GenerateAddress()
			}

			_, _, err = logic.Send(SenderAccount, rcvAddr, tc.transferAmount, tc.tipAmount, tc.gasLimit, tc.gasPrice, tc.contract, bc)

			assert.Equal(t, tc.expectedErr, err)

			//a short delay before mining starts
			time.Sleep(time.Millisecond * 500)

			// Make logic.Sender the miner and mine for 1 block (which should include the transaction)
			pow.Start()
			for bc.GetMaxHeight() < 1 {
			}
			pow.Stop()
			util.WaitDoneOrTimeout(func() bool {
				return !pow.IsProducingBlock()
			}, 20)
			// Verify balance of logic.Sender's account (genesis "mineReward" - transferred amount)
			SenderBalance, err := logic.GetBalance(SenderAccount.GetKeyPair().GenerateAddress(), bc)
			if err != nil {
				panic(err)
			}
			expectedBalance, _ := mineReward.Sub(tc.expectedTransfer)
			expectedBalance, _ = expectedBalance.Sub(tc.expectedTip)
			assert.Equal(t, expectedBalance, SenderBalance)

			// Balance of the miner's account should be the amount tipped + mineReward
			minerBalance, err := logic.GetBalance(minerAccount.GetKeyPair().GenerateAddress(), bc)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, mineReward.Times(bc.GetMaxHeight()).Add(tc.expectedTip), minerBalance)

			//check smart contract deployment
			res := string("")
			contractAddr := account.NewAddress("")
		loop:
			for i := bc.GetMaxHeight(); i > 0; i-- {
				blk, err := bc.GetBlockByHeight(i)
				assert.Nil(t, err)
				for _, tx := range blk.GetTransactions() {
					contractAddr = tx.GetContractAddress()
					if contractAddr.String() != "" {
						res = tx.Vout[transaction.ContractTxouputIndex].Contract
						break loop
					}
				}
			}
			assert.Equal(t, tc.contract, res)

			// Balance of the receiver's account should be the amount transferred
			var receiverBalance *common.Amount
			if isContract {
				receiverBalance, err = logic.GetBalance(contractAddr, bc)
			} else {
				receiverBalance, err = logic.GetBalance(receiverAccount.GetKeyPair().GenerateAddress(), bc)
			}
			assert.Equal(t, tc.expectedTransfer, receiverBalance)
		})
	}
}

//test logic.Send to invalid address
func TestSendToInvalidAddress(t *testing.T) {
	cleanUpDatabase()

	store := storage.NewRamStorage()
	defer store.Close()

	//this is internally set. Dont modify
	mineReward := common.NewAmount(10000000)
	//Transfer ammount
	transferAmount := common.NewAmount(25)
	tip := common.NewAmount(5)
	//create a account address
	account1, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
	assert.NotEmpty(t, account1)
	addr1 := account1.GetKeyPair().GenerateAddress()

	node := network.FakeNodeWithPidAndAddr(store, "test", "test")
	//create a blockchain
	bc, err := logic.CreateBlockchain(addr1, store, nil, transaction_pool.NewTransactionPool(node, 128), nil, 1000000)
	assert.Nil(t, err)
	assert.NotNil(t, bc)

	//The balance should be 10 after creating a blockchain
	balance1, err := logic.GetBalance(addr1, bc)
	assert.Nil(t, err)
	assert.Equal(t, mineReward, balance1)
	//pool := core.NewBlockPool()

	//bm := blockchain_logic.NewBlockchainManager(bc, pool, node)

	//logic.Send 5 coins from addr1 to an invalid address
	_, _, err = logic.Send(account1, account.NewAddress(InvalidAddress), transferAmount, tip, common.NewAmount(0), common.NewAmount(0), "", bc)

	assert.NotNil(t, err)

	//the balance of the first account should be still be 10
	balance1, err = logic.GetBalance(addr1, bc)
	assert.Nil(t, err)
	assert.Equal(t, mineReward, balance1)

	cleanUpDatabase()
}

//insufficient fund
func TestSendInsufficientBalance(t *testing.T) {
	cleanUpDatabase()

	store := storage.NewRamStorage()
	defer store.Close()

	tip := common.NewAmount(5)

	//this is internally set. Dont modify
	mineReward := common.NewAmount(10000000)
	//Transfer ammount is larger than the balance
	transferAmount := common.NewAmount(250000000)

	//create a account address
	account1, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
	assert.NotEmpty(t, account1)
	addr1 := account1.GetKeyPair().GenerateAddress()

	node := network.FakeNodeWithPidAndAddr(store, "test", "test")

	//create a blockchain
	bc, err := logic.CreateBlockchain(addr1, store, nil, transaction_pool.NewTransactionPool(node, 128), nil, 1000000)
	assert.Nil(t, err)
	assert.NotNil(t, bc)

	//The balance should be 10 after creating a blockchain
	balance1, err := logic.GetBalance(addr1, bc)
	assert.Nil(t, err)
	assert.Equal(t, mineReward, balance1)

	//Create a second account
	account2, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
	assert.NotEmpty(t, account2)
	assert.Nil(t, err)
	addr2 := account2.GetKeyPair().GenerateAddress()

	//The balance should be 0
	balance2, err := logic.GetBalance(addr2, bc)
	assert.Nil(t, err)
	assert.Equal(t, common.NewAmount(0), balance2)
	//pool := core.NewBlockPool()

	//bm := blockchain_logic.NewBlockchainManager(bc, pool, node)

	//logic.Send 5 coins from addr1 to addr2
	_, _, err = logic.Send(account1, addr2, transferAmount, tip, common.NewAmount(0), common.NewAmount(0), "", bc)

	assert.NotNil(t, err)

	//the balance of the first account should be still be 10
	balance1, err = logic.GetBalance(addr1, bc)
	assert.Nil(t, err)
	assert.Equal(t, mineReward, balance1)

	//the balance of the second account should be 0
	balance2, err = logic.GetBalance(addr2, bc)
	assert.Nil(t, err)
	assert.Equal(t, common.NewAmount(0), balance2)

	cleanUpDatabase()
}

func TestBlockMsgRelaySingleMiner(t *testing.T) {
	const (
		timeBetweenBlock = 1
	)
	cleanUpDatabase()
	var dposArray []*consensus.DPOS
	var bcs []*blockchain_logic.Blockchain
	var nodes []*network.Node

	validProducerAddr := "dastXXWLe5pxbRYFhcyUq8T3wb5srWkHKa"
	validProducerKey := "300c0338c4b0d49edc66113e3584e04c6b907f9ded711d396d522aae6a79be1a"

	producerAddrs := []string{}
	producerKey := []string{}
	numOfNodes := 4
	for i := 0; i < numOfNodes; i++ {
		producerAddrs = append(producerAddrs, validProducerAddr)
		producerKey = append(producerKey, validProducerKey)
	}

	dynasty := consensus.NewDynasty(producerAddrs, numOfNodes, timeBetweenBlock)

	for i := 0; i < numOfNodes; i++ {
		dpos := consensus.NewDPOS(block_producer_info.NewBlockProducerInfo(producerAddrs[0]))
		dpos.SetDynasty(dynasty)

		db := storage.NewRamStorage()
		node := network.NewNode(db, nil)
		node.Start(testport_msg_relay_port1+i, "")

		bc := blockchain_logic.CreateBlockchain(account.NewAddress(producerAddrs[0]), db, dpos, transaction_pool.NewTransactionPool(node, 128), nil, 100000)
		bcs = append(bcs, bc)

		nodes = append(nodes, node)

		dpos.SetKey(producerKey[0])
		dposArray = append(dposArray, dpos)
	}
	//each node connects to the subsequent node only
	for i := 0; i < len(nodes)-1; i++ {
		connectNodes(nodes[i], nodes[i+1])
	}

	//firstNode Starts Mining
	dposArray[0].Start()
	util.WaitDoneOrTimeout(func() bool {
		return bcs[0].GetMaxHeight() >= 5
	}, 8)
	dposArray[0].Stop()

	//expect every node should have # of entries in dapmsg cache equal to their blockchain height
	for i := 0; i < len(nodes)-1; i++ {
		assert.Equal(t, bcs[i].GetMaxHeight(), bcs[i+1].GetMaxHeight())
	}

	for _, node := range nodes {
		node.Stop()
	}
	cleanUpDatabase()
}

// Test if network radiation bounces forever
func TestBlockMsgRelayMeshNetworkMultipleMiners(t *testing.T) {
	const (
		timeBetweenBlock = 1
		dposRounds       = 2
		bufferTime       = 0
	)
	cleanUpDatabase()
	var dposArray []*consensus.DPOS
	var bcs []*blockchain_logic.Blockchain
	var nodes []*network.Node

	var firstNode *network.Node

	validProducerAddr := "dastXXWLe5pxbRYFhcyUq8T3wb5srWkHKa"
	validProducerKey := "300c0338c4b0d49edc66113e3584e04c6b907f9ded711d396d522aae6a79be1a"

	producerAddrs := []string{}
	producerKey := []string{}
	numOfNodes := 4
	for i := 0; i < numOfNodes; i++ {
		producerAddrs = append(producerAddrs, validProducerAddr)
		producerKey = append(producerKey, validProducerKey)
	}

	dynasty := consensus.NewDynasty(producerAddrs, numOfNodes, timeBetweenBlock)

	for i := 0; i < numOfNodes; i++ {
		dpos := consensus.NewDPOS(block_producer_info.NewBlockProducerInfo(producerAddrs[0]))
		dpos.SetDynasty(dynasty)

		db := storage.NewRamStorage()
		node := network.NewNode(db, nil)
		node.Start(testport_msg_relay_port2+i, "")
		bc := blockchain_logic.CreateBlockchain(account.NewAddress(producerAddrs[0]), storage.NewRamStorage(), dpos, transaction_pool.NewTransactionPool(node, 128), nil, 100000)
		bcs = append(bcs, bc)

		if i == 0 {
			firstNode = node
		} else {
			node.GetNetwork().ConnectToSeed(firstNode.GetHostPeerInfo())
		}
		nodes = append(nodes, node)

		dpos.SetKey(producerKey[0])
		dposArray = append(dposArray, dpos)
	}

	//each node connects to every other node
	for i := range nodes {
		for j := range nodes {
			if i != j {
				connectNodes(nodes[i], nodes[j])
			}
		}
	}

	//firstNode Starts Mining
	for _, dpos := range dposArray {
		dpos.Start()
	}

	time.Sleep(time.Second * time.Duration(dynasty.GetDynastyTime()*dposRounds+bufferTime))

	for _, dpos := range dposArray {
		dpos.Stop()
	}
	//expect every node should have # of entries in dapmsg cache equal to their blockchain height
	for i := 0; i < len(nodes)-1; i++ {
		assert.Equal(t, bcs[i].GetMaxHeight(), bcs[i+1].GetMaxHeight())
	}

	for _, node := range nodes {
		node.Stop()
	}

	cleanUpDatabase()
}

func TestForkChoice(t *testing.T) {
	var pows []*consensus.ProofOfWork
	var bcs []*blockchain_logic.Blockchain
	var bms []*blockchain_logic.BlockchainManager
	var dbs []storage.Storage
	var pools []*core.BlockPool
	// Remember to close all opened databases after test
	defer func() {
		for _, db := range dbs {
			db.Close()
		}
	}()

	addr := account.NewAddress("17DgRtQVvaytkiKAfXx9XbV23MESASSwUz")
	//wait for mining for at least "targetHeight" blocks
	//targetHeight := uint64(4)
	//num of nodes to be created in the test
	numOfNodes := 2
	var nodes []*network.Node
	for i := 0; i < numOfNodes; i++ {
		db := storage.NewRamStorage()
		dbs = append(dbs, db)

		node := network.NewNode(db, nil)
		bc, pow := CreateBlockchain(addr, addr, db, transaction_pool.NewTransactionPool(node, 128))
		bcs = append(bcs, bc)
		pool := core.NewBlockPool()
		pools = append(pools, pool)

		bm := blockchain_logic.NewBlockchainManager(bc, pool, node)

		pow.SetTargetBit(10)
		node.Start(testport_fork+i, "")
		pows = append(pows, pow)
		nodes = append(nodes, node)
		bms = append(bms, bm)
	}
	defer nodes[0].Stop()
	defer nodes[1].Stop()

	// Mine more blocks on node[0] than on node[1]
	pows[1].Start()
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetMaxHeight() > 4
	}, 10)
	pows[1].Stop()
	desiredHeight := uint64(10)
	if bcs[1].GetMaxHeight() > 10 {
		desiredHeight = bcs[1].GetMaxHeight() + 1
	}
	pows[0].Start()
	util.WaitDoneOrTimeout(func() bool {
		return bcs[0].GetMaxHeight() > desiredHeight
	}, 20)
	pows[0].Stop()

	util.WaitDoneOrTimeout(func() bool {
		return !pows[0].IsProducingBlock()
	}, 5)

	// Trigger fork choice in node[1] by broadcasting tail block of node[0]
	tailBlk, _ := bcs[0].GetTailBlock()
	connectNodes(nodes[0], nodes[1])
	bms[0].BroadcastBlock(tailBlk)
	// Make sure syncing starts on node[1]
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetState() == blockchain.BlockchainSync
	}, 10)
	// Make sure syncing ends on node[1]
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetState() != blockchain.BlockchainSync
	}, 20)

	assert.Equal(t, bcs[0].GetMaxHeight(), bcs[1].GetMaxHeight())
	assert.True(t, isSameBlockChain(bcs[0], bcs[1]))
}

func TestForkSegmentHandling(t *testing.T) {
	var pows []*consensus.ProofOfWork
	var bcs []*blockchain_logic.Blockchain
	var dbs []storage.Storage
	var pools []*core.BlockPool
	var bms []*blockchain_logic.BlockchainManager
	// Remember to close all opened databases after test
	defer func() {
		for _, db := range dbs {
			db.Close()
		}
	}()
	addr := account.NewAddress("17DgRtQVvaytkiKAfXx9XbV23MESASSwUz")

	numOfNodes := 2
	var nodes []*network.Node
	for i := 0; i < numOfNodes; i++ {
		db := storage.NewRamStorage()
		dbs = append(dbs, db)
		node := network.NewNode(db, nil)

		bc, pow := CreateBlockchain(addr, addr, db, transaction_pool.NewTransactionPool(node, 128))
		bcs = append(bcs, bc)
		pool := core.NewBlockPool()
		pools = append(pools, pool)
		bm := blockchain_logic.NewBlockchainManager(bc, pool, node)

		pow.SetTargetBit(10)
		node.Start(testport_fork_segment+i, "")
		pows = append(pows, pow)
		nodes = append(nodes, node)
		bms = append(bms, bm)
	}
	defer nodes[0].Stop()
	defer nodes[1].Stop()

	blk1 := &block.Block{}
	blk2 := &block.Block{}

	// Ensure node[1] mined some blocks
	pows[1].Start()
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetMaxHeight() > 3
	}, 10)
	pows[1].Stop()

	// Ensure node[0] mines more blocks than node[1]
	pows[0].Start()
	util.WaitDoneOrTimeout(func() bool {
		return bcs[0].GetMaxHeight() > 12
	}, 30)
	pows[0].Stop()

	util.WaitDoneOrTimeout(func() bool {
		return !pows[0].IsProducingBlock()
	}, 5)

	// Pick 2 blocks from blockchain[0] which can trigger syncing on node[1]
	mid := uint64(7)
	if bcs[0].GetMaxHeight() < 7 {
		mid = bcs[0].GetMaxHeight() - 1
	}
	blk1, _ = bcs[0].GetBlockByHeight(mid)
	blk2, _ = bcs[0].GetTailBlock()

	connectNodes(nodes[0], nodes[1])
	bms[0].BroadcastBlock(blk1)
	// Wait for node[1] to start syncing
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetState() == blockchain.BlockchainSync
	}, 10)

	// node[0] broadcast higher block on the same fork and should trigger another sync on node[1]
	bms[0].BroadcastBlock(blk2)

	// Make sure previous syncing ends
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetState() != blockchain.BlockchainSync
	}, 10)
	// Make sure node[1] is syncing again
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetState() == blockchain.BlockchainSync
	}, 10)
	// Make sure syncing ends
	util.WaitDoneOrTimeout(func() bool {
		return bcs[1].GetState() != blockchain.BlockchainSync
	}, 10)

	assert.True(t, isSameBlockChain(bcs[0], bcs[1]))
}

// Integration test for adding balance
func TestAddBalance(t *testing.T) {
	testCases := []struct {
		name         string
		addAmount    *common.Amount
		expectedDiff *common.Amount
		expectedErr  error
	}{
		{"Add 5", common.NewAmount(5), common.NewAmount(5), nil},
		{"Add zero", common.NewAmount(0), common.NewAmount(0), logic.ErrInvalidAmount},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			store := storage.NewRamStorage()
			defer store.Close()

			// Create a coinbase address
			key := "bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa7e"
			minerKeyPair := account.GenerateKeyPairByPrivateKey(key)
			minerAccount := account.NewAccountByKey(minerKeyPair)

			addr := minerAccount.GetKeyPair().GenerateAddress()
			node := network.FakeNodeWithPidAndAddr(store, "a", "b")

			bc, pow := CreateBlockchain(addr, addr, store, transaction_pool.NewTransactionPool(node, 128))

			// Create a new account address for testing
			testAddr := account.NewAddress("dGDrVKjCG3sdXtDUgWZ7Fp3Q97tLhqWivf")

			logic.SetMinerKeyPair(key)
			pow.SetTargetBit(0)
			pow.Start()

			for bc.GetMaxHeight() <= 1 {
			}

			// Add `addAmount` to the balance of the new account
			_, _, err := logic.SendFromMiner(testAddr, tc.addAmount, bc)
			height := bc.GetMaxHeight()
			assert.Equal(t, err, tc.expectedErr)
			for bc.GetMaxHeight()-height <= 1 {
			}

			pow.Stop()

			// The account balance should be the expected difference
			balance, err := logic.GetBalance(testAddr, bc)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedDiff, balance)
		})
	}
}

// Integration test for adding balance to invalid address
func TestAddBalanceWithInvalidAddress(t *testing.T) {
	testCases := []struct {
		name    string
		address string
	}{
		{"Invalid char in address", InvalidAddress},
		{"Invalid checksum address", "1AUrNJCRM5X5fDdmm3E3yjCrXQMLwfwfww"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := storage.NewRamStorage()
			defer db.Close()

			// Create a coinbase wallet address
			addr := account.NewAddress("dG6HhzSdA5m7KqvJNszVSf8i5f4neAteSs")
			node := network.FakeNodeWithPidAndAddr(db, "a", "b")
			// Create a blockchain
			bc, err := logic.CreateBlockchain(addr, db, nil, transaction_pool.NewTransactionPool(node, 128), nil, 1000000)
			assert.Nil(t, err)

			_, _, err = logic.SendFromMiner(account.NewAddress(tc.address), common.NewAmount(8), bc)
			assert.Equal(t, logic.ErrInvalidRcverAddress, err)
		})
	}
}

func TestSmartContractLocalStorage(t *testing.T) {
	store := storage.NewRamStorage()
	defer store.Close()

	contract := `'use strict';

	var StorageTest = function(){

	};

	StorageTest.prototype = {
	set:function(key,value){
			return LocalStorage.set(key,value);
		},
	get:function(key){
			return LocalStorage.get(key);
		}
	};
	var storageTest = new StorageTest;
	`

	// Create a account address
	SenderAccount, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
	minerAccount, err := logic.CreateAccount(logic.GetTestAccountPath(), "test")
	assert.Nil(t, err)
	node := network.FakeNodeWithPidAndAddr(store, "test", "test")
	bc, pow := CreateBlockchain(minerAccount.GetKeyPair().GenerateAddress(), SenderAccount.GetKeyPair().GenerateAddress(), store, transaction_pool.NewTransactionPool(node, 128))

	//deploy smart contract
	_, _, err = logic.Send(SenderAccount, account.NewAddress(""), common.NewAmount(1), common.NewAmount(0), common.NewAmount(10000), common.NewAmount(1), contract, bc)

	assert.Nil(t, err)

	txp := bc.GetTxPool().GetTransactions()[0]
	contractAddr := txp.GetContractAddress()

	// Create a miner account; Balance is 0 initially

	if err != nil {
		panic(err)
	}

	//a short delay before mining starts
	time.Sleep(time.Millisecond * 500)

	// Make logic.Sender the miner and mine for 1 block (which should include the transaction)
	pow.Start()
	for bc.GetMaxHeight() < 1 {
	}
	pow.Stop()

	//a short delay before mining starts
	time.Sleep(time.Millisecond * 500)

	//store data
	functionCall := `{"function":"set","args":["testKey","222"]}`
	_, _, err = logic.Send(SenderAccount, contractAddr, common.NewAmount(1), common.NewAmount(0), common.NewAmount(100), common.NewAmount(1), functionCall, bc)

	assert.Nil(t, err)
	pow.Start()
	for bc.GetMaxHeight() < 1 {
	}
	pow.Stop()

	//get data
	functionCall = `{"function":"get","args":["testKey"]}`
	_, _, err = logic.Send(SenderAccount, contractAddr, common.NewAmount(1), common.NewAmount(0), common.NewAmount(100), common.NewAmount(1), functionCall, bc)

	assert.Nil(t, err)
	pow.Start()
	for bc.GetMaxHeight() < 1 {
	}
	pow.Stop()
}

func connectNodes(node1 *network.Node, node2 *network.Node) {
	node1.GetNetwork().ConnectToSeed(node2.GetHostPeerInfo())
}

func setupNode(addr account.Address, pow *consensus.ProofOfWork, bc *blockchain_logic.Blockchain, port int) *network.Node {
	node := network.NewNode(bc.GetDb(), nil)

	pow.SetTargetBit(12)
	node.Start(port, "")
	defer node.Stop()
	return node
}

func CreateBlockchain(producer, addr account.Address, db *storage.RamStorage, txPool *transaction_pool.TransactionPool) (*blockchain_logic.Blockchain, *consensus.ProofOfWork) {
	pow := consensus.NewProofOfWork(block_producer_info.NewBlockProducerInfo(producer.String()))
	return blockchain_logic.CreateBlockchain(addr, db, pow, txPool, nil, 100000), pow
}

func TestDoubleMint(t *testing.T) {
	var SendNode *network.Node
	var recvNode *network.Node
	var recvNodeBc *blockchain_logic.Blockchain
	var blks []*block.Block
	var parent *block.Block
	var dposArray []*consensus.DPOS
	var SendBm *blockchain_logic.BlockchainManager

	validProducerAddr := "dPGZmHd73UpZhrM6uvgnzu49ttbLp4AzU8"
	validProducerKey := "5a66b0fdb69c99935783059bb200e86e97b506ae443a62febd7d0750cd7fac55"

	dynasty := consensus.NewDynasty([]string{validProducerAddr}, len([]string{validProducerAddr}), 15)
	producerHash, _ := account.GeneratePubKeyHashByAddress(account.NewAddress(validProducerAddr))
	tx := &transaction.Transaction{nil, []transaction_base.TXInput{{[]byte{}, -1, nil, nil}}, []transaction_base.TXOutput{{common.NewAmount(0), account.PubKeyHash(producerHash), ""}}, common.NewAmount(0), common.NewAmount(0), common.NewAmount(0)}

	for i := 0; i < 3; i++ {
		blk := createValidBlock([]*transaction.Transaction{tx}, validProducerKey, validProducerAddr, parent)
		blks = append(blks, blk)
		parent = blk
	}
	//check all timestamps are equal
	for i := 0; i < len(blks)-1; i++ {
		assert.True(t, blks[i].GetTimestamp() == blks[i+1].GetTimestamp())
	}
	for i := 0; i < 2; i++ {

		dpos := consensus.NewDPOS(block_producer_info.NewBlockProducerInfo(validProducerAddr))
		dpos.SetDynasty(dynasty)

		db := storage.NewRamStorage()
		node := network.NewNode(db, nil)
		node.Start(testport_msg_relay_port3+i, "")

		bc := blockchain_logic.CreateBlockchain(account.NewAddress(validProducerAddr), db, dpos, transaction_pool.NewTransactionPool(node, 128), nil, 100000)
		pool := core.NewBlockPool()

		bm := blockchain_logic.NewBlockchainManager(bc, pool, node)

		dpos.SetKey(validProducerKey)
		if i == 0 {
			SendNode = node
			SendBm = bm
		} else {
			recvNode = node
			recvNode.GetNetwork().ConnectToSeed(SendNode.GetHostPeerInfo())
			recvNodeBc = bc
		}
		dposArray = append(dposArray, dpos)
	}

	defer recvNode.Stop()
	defer SendNode.Stop()

	for _, blk := range blks {
		SendBm.BroadcastBlock(blk)
	}

	time.Sleep(time.Second * 2)
	assert.True(t, recvNodeBc.GetMaxHeight() < 2)
	assert.False(t, dposArray[1].Validate(blks[1]))
}

func createValidBlock(tx []*transaction.Transaction, validProducerKey, validProducerAddr string, parent *block.Block) *block.Block {
	blk := block.NewBlock(tx, parent, validProducerAddr)
	blk.SetHash(block_logic.CalculateHashWithNonce(blk))
	block_logic.SignBlock(blk, validProducerKey)
	return blk
}

func TestSimultaneousSyncingAndBlockProducing(t *testing.T) {
	const genesisAddr = "121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD"

	validProducerAddress := "dPGZmHd73UpZhrM6uvgnzu49ttbLp4AzU8"
	validProducerKey := "5a66b0fdb69c99935783059bb200e86e97b506ae443a62febd7d0750cd7fac55"
	//fmt.Println(validProducerAddress, validProducerKey)
	conss := consensus.NewDPOS(block_producer_info.NewBlockProducerInfo(validProducerAddress))
	dynasty := consensus.NewDynasty([]string{validProducerAddress}, 1, 1)
	conss.SetDynasty(dynasty)

	db := storage.NewRamStorage()
	seedNode := network.NewNode(db, nil)
	seedNode.Start(testport_fork_syncing, "")
	defer seedNode.Stop()

	bc := blockchain_logic.CreateBlockchain(account.NewAddress(genesisAddr), storage.NewRamStorage(), conss, transaction_pool.NewTransactionPool(seedNode, 128), nil, 100000)

	//create and start seed node
	pool := core.NewBlockPool()
	bm := blockchain_logic.NewBlockchainManager(bc, pool, seedNode)

	conss.SetKey(validProducerKey)

	// seed node start mining
	conss.Start()
	util.WaitDoneOrTimeout(func() bool {
		return bc.GetMaxHeight() > 8
	}, 10)

	// set up another node for syncing
	dpos := consensus.NewDPOS(block_producer_info.NewBlockProducerInfo(validProducerAddress))
	dpos.SetDynasty(dynasty)

	db1 := storage.NewRamStorage()
	node := network.NewNode(db1, nil)
	node.Start(testport_fork_syncing+1, "")
	defer node.Stop()

	bc1 := blockchain_logic.CreateBlockchain(account.NewAddress(genesisAddr), db1, dpos, transaction_pool.NewTransactionPool(node, 128), nil, 100000)

	dpos.SetKey(validProducerKey)

	// Trigger fork choice in node by broadcasting tail block of node[0]
	tailBlk, _ := bc.GetTailBlock()

	connectNodes(seedNode, node)
	bm.BroadcastBlock(tailBlk)

	time.Sleep(time.Second * 5)
	conss.Stop()
	assert.True(t, bc.GetMaxHeight()-bc1.GetMaxHeight() <= 1)
}

func TestUpdate(t *testing.T) {
	var address1Bytes = []byte("address1000000000000000000000000")
	var address1Hash, _ = account.NewUserPubKeyHash(address1Bytes)

	db := storage.NewRamStorage()
	defer db.Close()

	blk := core.GenerateUtxoMockBlockWithoutInputs()
	utxoIndex := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(db))
	utxoIndex.UpdateUtxoState(blk.GetTransactions())
	utxoIndex.Save()
	utxoIndexInDB := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(db))

	// test updating UTXO index with non-dependent transactions
	// Assert that both the original instance and the database copy are updated correctly
	for _, index := range []utxo_logic.UTXOIndex{*utxoIndex, *utxoIndexInDB} {
		utxoTx := index.GetAllUTXOsByPubKeyHash(address1Hash)
		assert.Equal(t, 2, utxoTx.Size())
		utxo0 := utxoTx.GetUtxo(blk.GetTransactions()[0].ID, 0)
		utx1 := utxoTx.GetUtxo(blk.GetTransactions()[0].ID, 1)
		assert.Equal(t, blk.GetTransactions()[0].ID, utxo0.Txid)
		assert.Equal(t, 0, utxo0.TxIndex)
		assert.Equal(t, blk.GetTransactions()[0].Vout[0].Value, utxo0.Value)
		assert.Equal(t, blk.GetTransactions()[0].ID, utx1.Txid)
		assert.Equal(t, 1, utx1.TxIndex)
		assert.Equal(t, blk.GetTransactions()[0].Vout[1].Value, utx1.Value)
	}

	// test updating UTXO index with dependent transactions
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

	var dependentTx1 = transaction.Transaction{
		ID: nil,
		Vin: []transaction_base.TXInput{
			{util.GenerateRandomAoB(1), 1, nil, pubkey1},
		},
		Vout: []transaction_base.TXOutput{
			{common.NewAmount(5), pkHash1, ""},
			{common.NewAmount(10), pkHash2, ""},
		},
		Tip: common.NewAmount(3),
	}
	dependentTx1.ID = dependentTx1.Hash()

	var dependentTx2 = transaction.Transaction{
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

	var dependentTx3 = transaction.Transaction{
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

	var dependentTx4 = transaction.Transaction{
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

	var dependentTx5 = transaction.Transaction{
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

	utxoPk2 := &utxo.UTXO{dependentTx1.Vout[1], dependentTx1.ID, 1, utxo.UtxoNormal}
	utxoPk1 := &utxo.UTXO{dependentTx1.Vout[0], dependentTx1.ID, 0, utxo.UtxoNormal}

	utxoTxPk2 := utxo.NewUTXOTx()
	utxoTxPk2.PutUtxo(utxoPk2)

	utxoTxPk1 := utxo.NewUTXOTx()
	utxoTxPk1.PutUtxo(utxoPk1)

	utxoIndex2 := utxo_logic.NewUTXOIndex(utxo.NewUTXOCache(storage.NewRamStorage()))

	utxoIndex2.SetIndex(map[string]*utxo.UTXOTx{
		pkHash2.String(): &utxoTxPk2,
		pkHash1.String(): &utxoTxPk1,
	})

	tx2Utxo1 := utxo.UTXO{dependentTx2.Vout[0], dependentTx2.ID, 0, utxo.UtxoNormal}
	tx2Utxo2 := utxo.UTXO{dependentTx2.Vout[1], dependentTx2.ID, 1, utxo.UtxoNormal}
	tx2Utxo3 := utxo.UTXO{dependentTx3.Vout[0], dependentTx3.ID, 0, utxo.UtxoNormal}
	tx2Utxo4 := utxo.UTXO{dependentTx1.Vout[0], dependentTx1.ID, 0, utxo.UtxoNormal}
	tx2Utxo5 := utxo.UTXO{dependentTx4.Vout[0], dependentTx4.ID, 0, utxo.UtxoNormal}
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey2).GetPrivateKey(), utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash2).GetAllUtxos(), &dependentTx2)
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey3).GetPrivateKey(), []*utxo.UTXO{&tx2Utxo1}, &dependentTx3)
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey4).GetPrivateKey(), []*utxo.UTXO{&tx2Utxo2, &tx2Utxo3}, &dependentTx4)
	transaction_logic.Sign(account.GenerateKeyPairByPrivateKey(prikey1).GetPrivateKey(), []*utxo.UTXO{&tx2Utxo4, &tx2Utxo5}, &dependentTx5)

	txsForUpdate := []*transaction.Transaction{&dependentTx2, &dependentTx3}
	utxoIndex2.UpdateUtxoState(txsForUpdate)
	assert.Equal(t, 1, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash1).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash2).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash3).Size())
	assert.Equal(t, 2, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash4).Size())
	txsForUpdate = []*transaction.Transaction{&dependentTx2, &dependentTx3, &dependentTx4}
	utxoIndex2.UpdateUtxoState(txsForUpdate)
	assert.Equal(t, 2, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash1).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash2).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash3).Size())
	txsForUpdate = []*transaction.Transaction{&dependentTx2, &dependentTx3, &dependentTx4, &dependentTx5}
	utxoIndex2.UpdateUtxoState(txsForUpdate)
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash1).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash2).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash3).Size())
	assert.Equal(t, 0, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash4).Size())
	assert.Equal(t, 1, utxoIndex2.GetAllUTXOsByPubKeyHash(pkHash5).Size())
}

func Test_MultipleMinersWithDPOS(t *testing.T) {
	const (
		timeBetweenBlock = 2
		dposRounds       = 3
	)

	miners := []string{
		"dPGZmHd73UpZhrM6uvgnzu49ttbLp4AzU8",
		"dQEooMsqp23RkPsvZXj3XbsRh9BUyGz2S9",
		"dastXXWLe5pxbRYFhcyUq8T3wb5srWkHKa",
		"dUuPPYshbBgkzUrgScEHWvdGbSxC8z4R12",
		"dPGD4t6ibpmyKZnXH1TNbbPw98EDaaZq8C",
	}
	keystrs := []string{
		"5a66b0fdb69c99935783059bb200e86e97b506ae443a62febd7d0750cd7fac55",
		"bb23d2ff19f5b16955e8a24dca34dd520980fe3bddca2b3e1b56663f0ec1aa7e",
		"300c0338c4b0d49edc66113e3584e04c6b907f9ded711d396d522aae6a79be1a",
		"da9282440fae188c371165e01615a2e1b14af68b3eaae51e6608c0bd86d4e6a6",
		"7c918ed7660d55759b7fc42b25f26bdab3caf8fc07586b2659a26470fb8dfc69",
	}
	dynasty := consensus.NewDynasty(miners, len(miners), timeBetweenBlock)
	var bps []*block_producer.BlockProducer
	var nodeArray []*network.Node

	for i, miner := range miners {
		producer := block_producer_info.NewBlockProducerInfo(miner)
		dpos := consensus.NewDPOS(producer)
		dpos.SetKey(keystrs[i])
		dpos.SetDynasty(dynasty)
		bc := blockchain_logic.CreateBlockchain(account.NewAddress(miners[0]), storage.NewRamStorage(), dpos, transaction_pool.NewTransactionPool(nil, 128), nil, 100000)
		pool := core.NewBlockPool()

		node := network.NewNode(bc.GetDb(), nil)
		node.Start(21200+i, "")
		nodeArray = append(nodeArray, node)

		bm := blockchain_logic.NewBlockchainManager(bc, pool, node)
		bp := block_producer.NewBlockProducer(bm, dpos, producer)
		bp.Start()
		bps = append(bps, bp)
	}

	for i := range miners {
		for j := range miners {
			if i != j {
				nodeArray[i].GetNetwork().ConnectToSeed(nodeArray[j].GetHostPeerInfo())
			}
		}
		bps[i].Start()
	}

	time.Sleep(time.Second*time.Duration(dynasty.GetDynastyTime()*dposRounds) + time.Second/2)

	for i := range miners {
		bps[i].Stop()
		nodeArray[i].Stop()
	}

	//Waiting block sync to other nodes
	time.Sleep(time.Second * 2)
	for i := range miners {
		v := bps[i]
		util.WaitDoneOrTimeout(func() bool {
			return !v.IsProducingBlock()
		}, 20)
	}

	for i := range miners {
		assert.Equal(t, uint64(dynasty.GetDynastyTime()*dposRounds/timeBetweenBlock), bps[i].Getblockchain().GetMaxHeight())
	}
}

func cleanUpDatabase() {
	account_logic.RemoveAccountFile()
}

func isSameBlockChain(bc1, bc2 *blockchain_logic.Blockchain) bool {
	if bc1 == nil || bc2 == nil {
		return false
	}

	bci1 := bc1.Iterator()
	bci2 := bc2.Iterator()
	if bc1.GetMaxHeight() != bc2.GetMaxHeight() {
		return false
	}

loop:
	for {
		blk1, _ := bci1.Next()
		blk2, _ := bci2.Next()
		if blk1 == nil || blk2 == nil {
			break loop
		}
		if !reflect.DeepEqual(blk1.GetHash(), blk2.GetHash()) {
			return false
		}
	}
	return true
}