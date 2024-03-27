package main

import (
	"context"
	"log"
	"net"

	"github.com/ziyw/simplekv/engine"
	pb "github.com/ziyw/simplekv/proto"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSimpleKeyValueServer
	bufferManager engine.BufferManager
}

func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	slog.Info("PUT request:",
		"ctx", ctx,
		"key", in.Key,
		"value", in.Value)
	err := s.bufferManager.Put(in.Key, in.Value)
	if err == nil {
		return &pb.PutResponse{Response: "DONE"}, nil
	}
	return nil, err
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	slog.Info("GET request:", "ctx", ctx, "key", in.Key)
	value, err := s.bufferManager.Get(in.Key)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{Value: value}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("error setup server, failed listen to: %v", err)
	}

	sv := grpc.NewServer()
	pb.RegisterSimpleKeyValueServer(sv, &Server{
		bufferManager: engine.BufferManager{},
	})
	slog.Info("server started", "addr:", lis.Addr())

	if err := sv.Serve(lis); err != nil {
		log.Fatalf("server crashed: %v", err)
	}
}
