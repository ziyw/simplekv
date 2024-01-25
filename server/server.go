package main

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"log"

	"net"

	pb "github.com/ziyw/simplekv/proto"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

const HASH_CONST = 10

type Server struct {
	pb.UnimplementedSimpleKeyValueServer
	activeSegment Segment
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()) % HASH_CONST
}

// Put (key,value) is to persist a key-value into a hashmap
// It will first hash the pair into a number by key, this hashnumber will be page_id
// read the page_id directly into a hashmap, insert key-value into the hashmap
// if hashmap is already hit the page size limit, return error
// then persist the hashmap on disk
func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	slog.Info("Receive Put (%v:%v)", in.Key, in.Value)

	s.activeSegment.Append(in.Key, in.Value)
	fmt.Println(s.activeSegment.hashmap)

	return &pb.PutResponse{Response: "DONE"}, nil

	// slog.Info(in.ProtoMessage())

	// page_id := strconv.Itoa(hash(in.Key))
	// page := Load(page_id)
	// fmt.Print("Existing Page Content is ")
	// fmt.Println(page.hashmap)

	// // TODO: hashmap fullness check
	// page.hashmap[in.Key] = in.Value
	// page.Flush()

	// fmt.Print("Updated page content is")
	// fmt.Println(page.hashmap)

	// // log.Printf("Receive Put Request (%v, %s)", in.GetKey(), in.GetValue())
	// return &pb.PutResponse{Response: "DONE"}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	slog.Info("Receive Get (%v)", in.Key)

	if val, ok := s.activeSegment.GetValue(in.Key); ok {
		return &pb.GetResponse{Value: val}, nil
	} else {
		return nil, errors.New("key doesn't not exist")
	}

}

// Start server
func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen to: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSimpleKeyValueServer(s, &Server{
		activeSegment: *NewSegment(1),
	})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
