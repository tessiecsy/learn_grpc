package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"BOOKSTORE_CLIENT/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 拨号 连接
	conn, err := grpc.Dial("127.0.0.1:8091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v\n", err)
	}
	defer conn.Close()

	// 创建客户端
	c := pb.NewBookstoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.ListBooks(ctx, &pb.ListBooksRequest{Shelf: 4})
	if err != nil {
		log.Fatalf("c.ListBooks failed, err:%v\n", err)
	}
	fmt.Printf("next_page_token:%v\n", res.NextPage)
	for i, book := range res.Books {
		fmt.Printf("%d: %#v\n", i, book)
	}

}
