package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/ziyw/simplekv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleKeyValueClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := client.Put(ctx, &pb.PutRequest{Key: "Hello", Value: "world"})
	if err != nil {
		log.Fatalf("client.Put failed: %v", err)
	}
	log.Printf("client.Put response: %v", r)

	v, err := client.Get(ctx, &pb.GetRequest{Key: "Hello"})
	if err != nil {
		log.Fatalf("client.Get response: %v", err)
	}
	log.Printf("%v", v)

}
