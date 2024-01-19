package main

import (
	"context"
	"log"
	"net"

	pb "github.com/ziyw/simplekv/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSimpleKeyValueServer
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Receive Get Request %v", in.GetKey())
	return &pb.GetResponse{Value: "value"}, nil
}

func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Receive Put Request (%v, %s)", in.GetKey(), in.GetValue())
	return &pb.PutResponse{Response: "DONE"}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen to: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSimpleKeyValueServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
