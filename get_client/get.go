package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/ziyw/simplekv/proto"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	k string
)

func init() {
	flag.StringVar(&k, "key", "", "key for DB query")
	flag.StringVar(&k, "k", "", "key for query")
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error dial to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleKeyValueClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Get(ctx, &pb.GetRequest{Key: k})
	slog.Info("Response: ", res, "Error:", err)
}
