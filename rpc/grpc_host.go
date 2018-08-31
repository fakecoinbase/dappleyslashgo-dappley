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
	"net"
	"log"
	"fmt"

	"golang.org/x/net/context"
	"github.com/dappley/go-dappley/rpc/pb"
	"github.com/dappley/go-dappley/network"
	"github.com/dappley/go-dappley/network/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	defaultRpcPort = 50051
	passwordToken  = "password"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct{
	srv 		*grpc.Server
	node 		*network.Node
	password	string
}

func NewGrpcServer(node *network.Node, adminPassword string) *Server{
	return &Server{grpc.NewServer(),node,adminPassword}
}

// SayHello implements helloworld.GreeterServer
func (s *Server) RpcCreateWallet(ctx context.Context, in *rpcpb.CreateWalletRequest) (*rpcpb.CreateWalletResponse, error) {
	return &rpcpb.CreateWalletResponse{Message: "Hello " + in.Name}, nil
}

func (s *Server) RpcGetBalance(ctx context.Context, in *rpcpb.GetBalanceRequest) (*rpcpb.GetBalanceResponse, error) {
	return &rpcpb.GetBalanceResponse{Message: "Hello " + in.Name}, nil
}

func (s *Server) RpcSend(ctx context.Context, in *rpcpb.SendRequest) (*rpcpb.SendResponse, error) {
	return &rpcpb.SendResponse{Message: "Hello " + in.Name}, nil
}

func (s *Server) RpcGetPeerInfo(ctx context.Context, in *rpcpb.GetPeerInfoRequest) (*rpcpb.GetPeerInfoResponse, error) {
	return &rpcpb.GetPeerInfoResponse{
		PeerList:s.node.GetPeerList().ToProto().(*networkpb.Peerlist),
	}, nil
}

func (s *Server) RpcGetBlockchainInfo(context.Context, *rpcpb.GetBlockchainInfoRequest) (*rpcpb.GetBlockchainInfoResponse, error){
	return &rpcpb.GetBlockchainInfoResponse{
		TailBlockHash: 	s.node.GetBlockchain().GetTailBlockHash(),
		BlockHeight:	s.node.GetBlockchain().GetMaxHeight(),
	}, nil
}

func (s *Server) Start(port uint32) {
	go func(){
		if port == 0{
			port = defaultRpcPort
		}
		lis, err := net.Listen("tcp", fmt.Sprint(":",port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		srv := grpc.NewServer(grpc.UnaryInterceptor(s.AuthInterceptor))
		rpcpb.RegisterConnectServer(srv, s)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
}

func (s *Server) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(meta[passwordToken]) != 1 {
		return nil, status.Errorf(codes.Unauthenticated, "No Password")
	}
	if meta[passwordToken][0] != s.password {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid Password")
	}
	return handler(ctx, req)
}

func (s *Server) Stop() {
	s.srv.Stop()
}
