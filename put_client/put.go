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
	v string
)

func init() {
	flag.StringVar(&k, "key", "", "key for DB query")
	flag.StringVar(&k, "k", "", "key for query")
	flag.StringVar(&v, "value", "", "value for query")
	flag.StringVar(&v, "v", "", "value for query")
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

	res, err := client.Put(ctx, &pb.PutRequest{Key: k, Value: v})
	slog.Info("Response: ", res, "Error:", err)
}
