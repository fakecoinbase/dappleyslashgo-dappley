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

package network

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	logger "github.com/sirupsen/logrus"

	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/core/pb"
	"github.com/dappley/go-dappley/network/pb"
)

const (
	protocalName           = "dappley/1.0.0"
	syncPeerTimeLimitMs    = 1000
	MaxMsgCountBeforeReset = 999999
	maxGetBlocksNum        = 10
)

var (
	ErrDapMsgNoCmd  = errors.New("command not specified")
	ErrIsInPeerlist = errors.New("peer already exists in peerlist")
)

type streamMsg struct {
	msg *DapMsg
	from peer.ID
}

type Node struct {
	host                   host.Host
	info                   *Peer
	bm                     *core.BlockChainManager
	streams                map[peer.ID]*Stream
	streamExitCh           chan *Stream
	peerList               *PeerList
	exitCh                 chan bool
	recentlyRcvedDapMsgs   *sync.Map
	dapMsgBroadcastCounter *uint64
	privKey                crypto.PrivKey
	dispatch				   chan *streamMsg
	downloadManager        *DownloadManager
}

func newMsg(dapMsg *DapMsg, id peer.ID) *streamMsg {
	return &streamMsg{dapMsg, id}
}

//create new Node instance
func NewNode(bc *core.Blockchain, pool *core.BlockPool) *Node {
	placeholder := uint64(0)
	bm := core.NewBlockChainManager()
	bm.SetblockPool(pool)
	bm.Setblockchain(bc)
	node := &Node{nil,
		nil,
		bm,
		make(map[peer.ID]*Stream, 10),
		make(chan *Stream, 10),
		NewPeerList(nil),
		make(chan bool, 1),
		&sync.Map{},
		&placeholder,
		nil,
		make(chan *streamMsg, 1000),
		nil,
	}
	node.downloadManager = NewDownloadManager(node)
	return node
}

func (n *Node) isNetworkRadiation(dapmsg DapMsg) bool {
	if _, value := n.recentlyRcvedDapMsgs.Load(dapmsg.GetKey()); value == true {
		return true
	}
	return false
}

func (n *Node) GetBlockchain() *core.Blockchain    { return n.bm.Getblockchain() }
func (n *Node) GetBlockPool() *core.BlockPool      { return n.bm.GetblockPool() }
func (n *Node) GetPeerList() *PeerList             { return n.peerList }
func (n *Node) GetRecentlyRcvedDapMsgs() *sync.Map { return n.recentlyRcvedDapMsgs }
func (n *Node) GetDownloadManager() *DownloadManager { return n.downloadManager }
func (n *Node) GetStreams() map[peer.ID]*Stream      { return n.streams }

func (n *Node) Start(listenPort int) error {

	h, addr, err := createBasicHost(listenPort, n.privKey)
	if err != nil {
		logger.WithError(err).Error("Node: failed to create basic host.")
		return err
	}

	n.host = h
	n.info, err = CreatePeerFromMultiaddr(addr)

	//set streamhandler. streamHanlder function is called upon stream connection
	n.host.SetStreamHandler(protocalName, n.streamHandler)
	n.StartRequestLoop()
	n.StartListenLoop()
	n.StartExitListener()
	return err
}

func (n *Node) Stop() {
	n.exitCh <- true
	for _, s := range n.streams {
		s.StopStream(nil)
	}
	n.host.RemoveStreamHandler(protocalName)
	err := n.host.Close()
	if err != nil {
		logger.WithError(err).Warn("Node: host was not closed properly.")
	}
}

func (n *Node) StartExitListener() {
	go func() {
		for {
			if s,ok := <-n.streamExitCh; ok {
				n.DisconnectPeer(s.peerID, s.remoteAddr)
			}
		}
	}()
}

func (n *Node) StartRequestLoop() {

	go func() {
		for {
			select {
			case <-n.exitCh:
				return
			case brPars := <-n.bm.GetblockPool().BlockRequestCh():
				n.RequestBlockUnicast(brPars.BlockHash, brPars.Pid)
			}
		}
	}()

}

func (n *Node) StartListenLoop() {
	go func() {
		for {
			if streamMsg, ok := <-n.dispatch; ok {
				n.handle(streamMsg.msg, streamMsg.from)
			}
		}
	}()

}

//LoadNetworkKeyFromFile reads the network privatekey from a file
func (n *Node) LoadNetworkKeyFromFile(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	data, err := base64.StdEncoding.DecodeString(string(bytes))
	if err != nil {
		return err
	}

	n.privKey, err = crypto.UnmarshalPrivateKey(data)
	if err != nil {
		return err
	}

	return nil
}

//create basic host. Returns host object, host address and error
func createBasicHost(listenPort int, priv crypto.PrivKey) (host.Host, ma.Multiaddr, error) {

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	}

	if priv != nil {
		opts = append(opts, libp2p.Identity(priv))
	}

	basicHost, err := libp2p.New(context.Background(), opts...)

	if err != nil {
		logger.WithError(err).Error("Failed to create a new libp2p node.")
		return nil, nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	logger.WithFields(logger.Fields{
		"address": fullAddr,
	}).Info("Host is up.")
	return basicHost, fullAddr, nil
}

//AddStreamsByString adds streams by their full addresses
func (n *Node) AddStreamsByString(targetFullAddrs []string){
	for _, fullAddr := range targetFullAddrs{
		err := n.AddStreamByString(fullAddr)
		if err!=nil{
			logger.WithError(err).WithFields(logger.Fields{
				"full_addr"	: fullAddr,
			}).Warn("Node: not able to add stream")
		}
	}
}

func (n *Node) AddStreamByString(targetFullAddr string) error {
	addr, err := ma.NewMultiaddr(targetFullAddr)
	if err != nil {
		return err
	}
	return n.AddStreamMultiAddr(addr)
}

//AddStreamMultiAddr stream to the targetFullAddr address. If the targetFullAddr is nil, the node goes to listening mode
func (n *Node) AddStreamMultiAddr(targetFullAddr ma.Multiaddr) error {

	//If there is a target address, connect to that address
	if targetFullAddr != nil {

		peerInfo, err := CreatePeerFromMultiaddr(targetFullAddr)
		if err != nil {
			return err
		}

		//Add Stream
		n.AddStream(peerInfo.peerid, peerInfo.addr)
	}

	return nil
}

func (n *Node) AddStream(peerid peer.ID, targetAddr ma.Multiaddr) error {
	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	p := Peer{peerid, targetAddr}
	if n.peerList.IsInPeerlist(&p) {
		logger.WithFields(logger.Fields{
			"host":   n.GetPeerMultiaddr().String(),
			"target": targetAddr.String(),
		}).Debug("Node: target has already been added before!")
		return ErrIsInPeerlist
	}

	n.host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	// make a new stream
	stream, err := n.host.NewStream(context.Background(), peerid, protocalName)
	if err != nil {
		return err
	}
	// Create a buffered stream so that read and write are non blocking.
	n.streamHandler(stream)

	// Add the peer list
	if n.peerList.ListIsFull() {
		n.peerList.RemoveRandomIP()
	}
	n.peerList.Add(&Peer{peerid, targetAddr})

	return nil
}

func (n *Node) DisconnectPeer(peerid peer.ID, targetAddr ma.Multiaddr) {
	delete(n.streams, peerid)
	n.peerList.DeletePeer(&Peer{peerid, targetAddr})
	if n.downloadManager != nil {
		n.downloadManager.DisconnectPeer(peerid)
	}
}

func (n *Node) handle(msg *DapMsg, id peer.ID) {
	switch msg.GetCmd() {
	case GetBlockchainInfo:
		n.GetBlockchainInfoHandler(msg, id)

	case ReturnBlockchainInfo:
		n.ReturnBlockchainInfoHandler(msg, id)

	case SyncBlock:
		n.SyncBlockHandler(msg, id)

	case SyncPeerList:
		n.AddMultiPeers(msg.GetData())

	case RequestBlock:
		n.SendRequestedBlock(msg.GetData(), id)
	case BroadcastTx:
		n.AddTxToPool(msg)

	case GetBlocks:
		n.GetBlocksHandler(msg, id)

	case ReturnBlocks:
		n.ReturnBlocksHandler(msg, id)

	default:
		logger.WithFields(logger.Fields{
			"from": id,
		}).Debug("Node: received an invalid command.")
	}
}

func (n *Node) streamHandler(s net.Stream) {
	// Create a buffer stream for non blocking read and write.
	peer := &Peer{s.Conn().RemotePeer(), s.Conn().RemoteMultiaddr()}
	if n.peerList.ListIsFull() {
		logger.WithFields(logger.Fields{
			"peers": len(n.peerList.GetPeerlist()),
		}).Warn("Node: peer list is full.")
	}
	if n.peerList.IsInPeerlist(peer) {
		logger.WithFields(logger.Fields{
			"peer": s.Conn().RemotePeer(),
		}).Warn("Node: peer is already in the peer list.")
	}
	logger.WithFields(logger.Fields{
		"host":   n.GetPeerID(),
		"target": s.Conn().RemotePeer(),
	}).Info("Node: is creating a new stream.")

	n.peerList.Add(peer)
	//start stream
	ns := NewStream(s)
	//add stream to this.streams
	n.streams[s.Conn().RemotePeer()] = ns
	ns.Start(n.streamExitCh, n.dispatch)

	n.SyncPeersUnicast(peer.peerid)
}

func (n *Node) GetInfo() *Peer { return n.info }

func (n *Node) GetPeerMultiaddr() ma.Multiaddr {
	if n.info == nil {
		return nil
	}
	return n.info.addr
}

func (n *Node) GetPeerID() peer.ID { return n.info.peerid }

func (n *Node) RelayDapMsg(dm DapMsg) {
	msgData := dm.ToProto()
	bytes, _ := proto.Marshal(msgData)
	n.broadcast(bytes)
}

func (n *Node) prepareData(msgData proto.Message, cmd string, uniOrBroadcast int, msgKey string) ([]byte, error) {
	if cmd == "" {
		return nil, ErrDapMsgNoCmd
	}

	bytes := []byte{}
	var err error

	if msgData != nil {
		//marshal the block to wire format
		bytes, err = proto.Marshal(msgData)
		if err != nil {
			return nil, err
		}
	}

	//build a dappley message
	dm := NewDapmsg(cmd, bytes, msgKey, uniOrBroadcast, n.dapMsgBroadcastCounter)
	if dm.cmd == SyncBlock || dm.cmd == BroadcastTx {
		n.cacheDapMsg(*dm)
	}
	data, err := proto.Marshal(dm.ToProto())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (n *Node) BroadcastBlock(block *core.Block) error {
	logger.WithFields(logger.Fields{
		"peer_id": n.GetPeerID(),
		"height":  block.GetHeight(),
		"hash":    hex.EncodeToString(block.GetHash()),
	}).Info("Node: is broadcasting a block.")
	data, err := n.prepareData(block.ToProto(), SyncBlock, Broadcast, hex.EncodeToString(block.GetHash()))
	if err != nil {
		return err
	}
	n.broadcast(data)
	return nil
}

func (n *Node) BroadcastGetBlockchainInfo() {
	request := &networkpb.GetBlockchainInfo{Version: protocalName}
	data, err := n.prepareData(request, GetBlockchainInfo, Broadcast, "")
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
		}).Warn("Broadcast GetBlockchainInfo failed.")
	}

	n.broadcast(data)
}

func (n *Node) SyncPeersBroadcast() error {
	data, err := n.prepareData(n.peerList.ToProto(), SyncPeerList, Broadcast, "")
	if err != nil {
		return err
	}
	n.broadcast(data)
	return nil
}

func (n *Node) TxBroadcast(tx *core.Transaction) error {
	data, err := n.prepareData(tx.ToProto(), BroadcastTx, Broadcast, hex.EncodeToString(tx.ID))
	if err != nil {
		return err
	}
	n.broadcast(data)
	return nil
}

func (n *Node) SyncPeersUnicast(pid peer.ID) error {
	data, err := n.prepareData(n.peerList.ToProto(), SyncPeerList, Unicast, "")
	if err != nil {
		return err
	}
	n.unicast(data, pid)
	return nil
}

func (n *Node) SendBlockUnicast(block *core.Block, pid peer.ID) error {
	data, err := n.prepareData(block.ToProto(), SyncBlock, Unicast, hex.EncodeToString(block.GetHash()))
	if err != nil {
		return err
	}
	n.unicast(data, pid)
	return nil
}

func (n *Node) RequestBlockUnicast(hash core.Hash, pid peer.ID) error {
	//build a deppley message

	dm := NewDapmsg(RequestBlock, hash, hex.EncodeToString(hash), Unicast, n.dapMsgBroadcastCounter)
	data, err := proto.Marshal(dm.ToProto())
	if err != nil {
		return err
	}
	n.unicast(data, pid)
	return nil
}

func (n *Node) DownloadBlocksUnicast(hashes []core.Hash, pid peer.ID) error {
	getBlockPb := &networkpb.GetBlocks{}

	getBlockPb.StartBlockHashes = make([][]byte, len(hashes))

	for index, hash := range hashes {
		getBlockPb.StartBlockHashes[index] = hash
	}

	data, err := n.prepareData(getBlockPb, GetBlocks, Unicast, "")
	if err != nil {
		return nil
	}

	n.unicast(data, pid)
	return nil
}

//broadcast data
func (n *Node) broadcast(data []byte) {
	for _, s := range n.streams {
		s.Send(data)
	}
}

//unicast data
func (n *Node) unicast(data []byte, pid peer.ID) {
	n.streams[pid].Send(data)
}

func (n *Node) addBlockToPool(block *core.Block, pid peer.ID) {
	//add block to blockpool. Make sure this is none blocking.
	n.bm.Push(block, pid)
}

func (n *Node) getFromProtoBlockMsg(data []byte) *core.Block {
	//create a block proto
	blockpb := &corepb.Block{}

	//unmarshal byte to proto
	if err := proto.Unmarshal(data, blockpb); err != nil {
		logger.Warn(err)
	}

	//create an empty block
	block := &core.Block{}

	//load the block with proto
	block.FromProto(blockpb)

	return block
}

func (n *Node) SyncBlockHandler(dm *DapMsg, pid peer.ID) {
	if len(dm.data) == 0 {
		logger.WithFields(logger.Fields{
			"cmd": "sync block",
		}).Warn("Node: cannot find block information.")
		return
	}

	if dm.uniOrBroadcast == Broadcast {
		if n.isNetworkRadiation(*dm) {
			return
		}
		n.cacheDapMsg(*dm)
		blk := n.getFromProtoBlockMsg(dm.GetData())
		n.addBlockToPool(blk, pid)
		if dm.uniOrBroadcast == Broadcast {
			n.RelayDapMsg(*dm)
		}
	} else {
		blk := n.getFromProtoBlockMsg(dm.GetData())
		n.addBlockToPool(blk, pid)
	}
}

func (n *Node) GetBlockchainInfoHandler(dm *DapMsg, pid peer.ID) {
	tailBlock, err := n.GetBlockchain().GetTailBlock()
	if err != nil {
		logger.WithFields(logger.Fields{
			"cmd": "GetBlockchainInfoRequest",
		}).Warn("Get tailblock failed.")
		return
	}

	response := &networkpb.ReturnBlockchainInfo{
		TailBlockHash: n.GetBlockchain().GetTailBlockHash(),
		BlockHeight:   n.GetBlockchain().GetMaxHeight(),
		Timestamp:     tailBlock.GetTimestamp(),
	}

	data, err := n.prepareData(response, ReturnBlockchainInfo, Unicast, "")
	if err != nil {
		logger.WithFields(logger.Fields{
			"cmd": "GetBlockchainInfoRequest",
		}).Warn("Prepare data failed.")
		return
	}

	n.unicast(data, pid)
}

func (n *Node) ReturnBlockchainInfoHandler(dm *DapMsg, pid peer.ID) {
	blockchainInfo := &networkpb.ReturnBlockchainInfo{}
	if err := proto.Unmarshal(dm.data, blockchainInfo); err != nil {
		logger.WithFields(logger.Fields{
			"cmd": "ReturnBlockchainInfo",
		}).Info("Parse data failed.")
		return
	}

	n.downloadManager.AddPeerBlockChainInfo(pid, blockchainInfo.BlockHeight)
}

//TODO  Refactor getblocks in rpcService and node
func (n *Node) GetBlocksHandler(dm *DapMsg, pid peer.ID) {
	param := &networkpb.GetBlocks{}
	if err := proto.Unmarshal(dm.data, param); err != nil {
		logger.WithFields(logger.Fields{
			"cmd": "GetBlocks",
		}).Info("Parse data failed.")
		return
	}

	block := n.findBlockInRequestHash(param.StartBlockHashes)

	// Reach the blockchain's tail
	if block.GetHeight() >= n.GetBlockchain().GetMaxHeight() {
		logger.WithFields(logger.Fields{
			"cmd": "GetBlocks",
		}).Info("Reach blockchain tail.")
		return
	}

	var blocks []*core.Block

	block, err := n.GetBlockchain().GetBlockByHeight(block.GetHeight() + 1)
	for i := int32(0); i < maxGetBlocksNum && err == nil; i++ {
		if block.GetHeight() == 0 {
			logger.Panicf("Error %v", hex.EncodeToString(block.GetHash()))
		}
		blocks = append(blocks, block)
		block, err = n.GetBlockchain().GetBlockByHeight(block.GetHeight() + 1)
	}

	result := &networkpb.ReturnBlocks{}
	for i := len(blocks) - 1; i >= 0; i-- {
		result.Blocks = append(result.Blocks, blocks[i].ToProto().(*corepb.Block))
	}

	result.StartBlockHashes = param.StartBlockHashes

	data, err := n.prepareData(result, ReturnBlocks, Unicast, "")
	if err != nil {
		logger.WithFields(logger.Fields{
			"cmd": "GetBlocks",
		}).Warn("Prepare data failed.")
		return
	}

	n.unicast(data, pid)
}

func (n *Node) findBlockInRequestHash(startBlockHashes [][]byte) *core.Block {
	for _, hash := range startBlockHashes {
		// hash in blockchain, return
		if block, err := n.GetBlockchain().GetBlockByHash(hash); err == nil {
			return block
		}
	}

	// Return Genesis Block
	block, _ := n.GetBlockchain().GetBlockByHeight(0)
	return block
}

func (n *Node) ReturnBlocksHandler(dm *DapMsg, pid peer.ID) {
	param := &networkpb.ReturnBlocks{}
	if err := proto.Unmarshal(dm.data, param); err != nil {
		logger.WithFields(logger.Fields{
			"cmd": "ReturnBlocks",
		}).Info("Parse data failed.")
		return
	}

	if n.downloadManager != nil {
		n.downloadManager.GetBlocksDataHandler(param, pid)
	}
}

func (n *Node) cacheDapMsg(dm DapMsg) {
	n.recentlyRcvedDapMsgs.Store(dm.GetKey(), true)
}

func (n *Node) AddTxToPool(dm *DapMsg) {
	if n.isNetworkRadiation(*dm) {
		return
	}

	n.RelayDapMsg(*dm)
	n.cacheDapMsg(*dm)

	//create a block proto
	txpb := &corepb.Transaction{}

	//unmarshal byte to proto
	if err := proto.Unmarshal(dm.GetData(), txpb); err != nil {
		logger.Warn(err)
	}

	//create an empty tx
	tx := &core.Transaction{}

	//load the tx with proto
	tx.FromProto(txpb)
	//add tx to txpool
	utxoIndex := core.LoadUTXOIndex(n.GetBlockchain().GetDb())
	utxoIndex.UpdateUtxoState(n.GetBlockchain().GetTxPool().GetTransactions())

	if tx.Verify(utxoIndex, 0) == false {
		logger.Info("Broadcast transaction verify failed.")
		return
	}
	n.bm.Getblockchain().GetTxPool().Push(tx)
}

func (n *Node) AddMultiPeers(data []byte) {

	go func() {
		//create a peerList proto
		plpb := &networkpb.Peerlist{}

		//unmarshal byte to proto
		if err := proto.Unmarshal(data, plpb); err != nil {
			logger.Warn(err)
		}

		//create an empty peerList
		pl := &PeerList{}

		//load the block with proto
		pl.FromProto(plpb)

		//remove the node's own peer info from the list
		newpl := &PeerList{[]*Peer{n.info}}
		newpl = newpl.FindNewPeers(pl)
		//find the new added peers
		newpl = n.peerList.FindNewPeers(newpl)

		//wait for random time within the time limit
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(syncPeerTimeLimitMs)))

		//add streams for new peers
		for _, p := range newpl.GetPeerlist() {
			if !n.peerList.IsInPeerlist(p) && p.peerid != n.info.peerid {
				n.AddStream(p.peerid, p.addr)
			}
		}

		//add peers
		n.peerList.MergePeerlist(newpl)
	}()
}

func (n *Node) SendRequestedBlock(hash []byte, pid peer.ID) {
	blockBytes, err := n.bm.Getblockchain().GetDb().Get(hash)
	if err != nil {
		logger.Warn("Node: failed to get the requested block from database.")
		return
	}
	block := core.Deserialize(blockBytes)
	n.SendBlockUnicast(block, pid)
}
