package main

import (
	"context"
	"errors"
	"log"

	"net"

	pb "github.com/ziyw/simplekv/proto"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

const HASH_CONST = 10

type Server struct {
	pb.UnimplementedSimpleKeyValueServer
	segList *SegList
}

func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	slog.Info("Put Request", "key", in.GetKey(), "value", in.GetValue())

	var seg *Segment
	if s.segList.Size() != 0 {
		seg = s.segList.GetCurrentSegment()
	} else {
		seg = NewSegment(s.segList.GetNextId())
		s.segList.AddSegment(seg)
	}

	seg.Append(in.Key, in.Value)

	slog.Info("Put Response", "response", "DONE")
	return &pb.PutResponse{Response: "DONE"}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	slog.Info("Get Request", "key", in.GetKey())

	if s.segList.Size() == 0 {
		return nil, errors.New("key doesn't not exist")
	}

	// traverse all segment to get value
	for i := s.segList.Size() - 1; i >= 0; i-- {
		if val, ok := s.segList.list[i].GetValue(in.Key); !ok {
			continue
		} else {
			slog.Info("Get Response", "value", val, "err", nil)
			return &pb.GetResponse{Value: val}, nil
		}
	}

	slog.Info("Get Response", "value", nil, "err", "NoExistingKey")
	return nil, errors.New("key doesn't not exist")
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen to: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSimpleKeyValueServer(s, &Server{
		segList: NewSegList(),
	})

	slog.Info("Started Server", "port", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
