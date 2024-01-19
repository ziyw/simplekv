package main

import (
	"context"
	"log"
	"time"

	pb "github.com/ziyw/simplekv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleKeyValueClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	v, err := client.Get(ctx, &pb.GetRequest{Key: "Hello"})
	if err != nil {
		log.Fatalf("client.Get failed: %v", err)
	}
	log.Printf("value is :%v", v)
}
