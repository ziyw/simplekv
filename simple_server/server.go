package main

import (
	"context"
	"log"
	"net"

	pb "github.com/ziyw/simplekv/proto"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSimpleKeyValueServer
}

func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	slog.Info("PUT request:", ctx, in.Key, in.Value)
	slog.Info("PUT response:DONE")
	return &pb.PutResponse{Response: "DONE"}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	slog.Info("GET request:", ctx, in.Key)
	slog.Info("GET response: DONE")
	return &pb.GetResponse{Value: "DONE"}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("error setup server, failed listen to: %v", err)
	}

	sv := grpc.NewServer()
	pb.RegisterSimpleKeyValueServer(sv, &Server{})
	slog.Info("server started", lis.Addr())

	if err := sv.Serve(lis); err != nil {
		log.Fatalf("server crashed: %v", err)
	}
}
