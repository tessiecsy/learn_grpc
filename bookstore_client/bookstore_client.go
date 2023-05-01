package main

import (
	"bookstore_client/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("grpc dial failed, err:%v", err)
		return
	}
	defer conn.Close()

	c := pb.NewBookstoreClient(conn)
	res, err := c.ListBooks(context.Background(), &pb.ListBooksRequest{Shelf: 4})
	if err != nil {
		log.Printf("c.ListBooks failed, err:%v", err)
		return
	}
	for i, book := range res.Books {
		log.Printf("%d:%v:", i, book)
	}

}
