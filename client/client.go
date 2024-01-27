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

var (
	f string
	k string
	v string
)

func init() {
	flag.StringVar(&f, "functions", "get", "get or put function")
	flag.StringVar(&f, "f", "get", "get or put function")

	flag.StringVar(&k, "key", "", "key for query")
	flag.StringVar(&k, "k", "", "key for query")

	flag.StringVar(&v, "value", "", "value for query")
	flag.StringVar(&v, "v", "", "value for query")
}

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

	switch f {
	case "get":
		{
			v, err := client.Get(ctx, &pb.GetRequest{Key: k})
			if err != nil {
				log.Fatalf("client.Get response: %v", err)
			}
			log.Printf("get response: %v", v)
			return
		}
	case "put":
		{
			v, err := client.Put(ctx, &pb.PutRequest{Key: k, Value: v})
			if err != nil {
				log.Fatalf("put failed: %v", err)
			}
			log.Printf("put response: %v", v)
		}
	default:
		log.Println("Function doesn't exist: " + f)
	}
}
