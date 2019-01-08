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
package rpc

import (
	"context"
	"encoding/hex"

	"github.com/dappley/go-dappley/client"
	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/logic"
	"github.com/dappley/go-dappley/network"
	"github.com/dappley/go-dappley/network/pb"
	"github.com/dappley/go-dappley/rpc/pb"
)

type AdminRpcService struct {
	node *network.Node
}

func (adminRpcService *AdminRpcService) RpcAddPeer(ctx context.Context, in *rpcpb.AddPeerRequest) (*rpcpb.AddPeerResponse, error) {
	status := "succeed"
	err := adminRpcService.node.AddStreamByString(in.FullAddress)
	if err != nil {
		status = err.Error()
	}
	return &rpcpb.AddPeerResponse{
		Status: status,
	}, nil
}

func (adminRpcService *AdminRpcService) RpcAddProducer(ctx context.Context, in *rpcpb.AddProducerRequest) (*rpcpb.AddProducerResponse, error) {
	if len(in.Address) == 0 {
		return &rpcpb.AddProducerResponse{
			Message: "address is empty",
		}, core.ErrInvalidAddress
	}
	err := adminRpcService.node.GetBlockchain().GetConsensus().AddProducer(in.Address)
	if err != nil {
		return &rpcpb.AddProducerResponse{
			Message: "failed to add producer: " + err.Error(),
		}, err
	}
	return &rpcpb.AddProducerResponse{
		Message: "producer is added",
	}, nil
}

func (adminRpcService *AdminRpcService) RpcGetPeerInfo(ctx context.Context, in *rpcpb.GetPeerInfoRequest) (*rpcpb.GetPeerInfoResponse, error) {
	return &rpcpb.GetPeerInfoResponse{
		PeerList: adminRpcService.node.GetPeerList().ToProto().(*networkpb.Peerlist),
	}, nil
}

//unlock the wallet through rpc service
func (adminRpcService *AdminRpcService) RpcUnlockWallet(ctx context.Context, in *rpcpb.UnlockWalletRequest) (*rpcpb.UnlockWalletResponse, error) {
	err := logic.SetUnLockWallet()
	if err != nil {
		return &rpcpb.UnlockWalletResponse{Message: err.Error()}, err
	}
	return &rpcpb.UnlockWalletResponse{Message: "succeed"}, nil
}

func (adminRpcService *AdminRpcService) RpcSendFromMiner(ctx context.Context, in *rpcpb.SendFromMinerRequest) (*rpcpb.SendFromMinerResponse, error) {
	sendToAddress := core.NewAddress(in.To)
	sendAmount := common.NewAmountFromBytes(in.Amount)
	if sendAmount.Validate() != nil || sendAmount.IsZero() {
		return &rpcpb.SendFromMinerResponse{Message: "invalid send amount (must be > 0)"}, core.ErrInvalidAmount
	}

	_, _, err := logic.SendFromMiner(sendToAddress, sendAmount, adminRpcService.node.GetBlockchain(), adminRpcService.node)
	if err != nil {
		return &rpcpb.SendFromMinerResponse{Message: "failed to add balance: " + err.Error()}, err
	}
	return &rpcpb.SendFromMinerResponse{Message: "succeed"}, nil
}

func (adminRpcService *AdminRpcService) RpcSend(ctx context.Context, in *rpcpb.SendRequest) (*rpcpb.SendResponse, error) {
	sendFromAddress := core.NewAddress(in.From)
	sendToAddress := core.NewAddress(in.To)
	sendAmount := common.NewAmountFromBytes(in.Amount)
	tip := common.NewAmountFromBytes(in.Tip)

	if sendAmount.Validate() != nil || sendAmount.IsZero() {
		return &rpcpb.SendResponse{Message: "invalid send amount (must be > 0)"}, core.ErrInvalidAmount
	}
	path := in.WalletPath
	if len(in.WalletPath) == 0 {
		path = client.GetWalletFilePath()
	}

	wm, err := logic.GetWalletManager(path)
	if err != nil {
		return &rpcpb.SendResponse{Message: "error loading local wallets"}, err
	}

	senderWallet := wm.GetWalletByAddress(sendFromAddress)
	if senderWallet == nil || len(senderWallet.Addresses) == 0 {
		return &rpcpb.SendResponse{Message: "sender wallet is not found"}, client.ErrAddressNotFound
	}

	txhash, scAddress, err := logic.Send(senderWallet, sendToAddress, sendAmount, tip, in.Data, adminRpcService.node.GetBlockchain(), adminRpcService.node)
	txhashStr := hex.EncodeToString(txhash)
	if err != nil {
		return &rpcpb.SendResponse{Message: "failed to send transaction", Txid: txhashStr}, err
	}

	resp := &rpcpb.SendResponse{Message: "succeed", Txid: txhashStr}
	if scAddress != "" {
		resp.ContractAddr = scAddress
	}

	return resp, nil
}

func (adminRpcService *AdminRpcService) IsPrivate() bool { return true }
