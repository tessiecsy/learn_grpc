package main

import (
	"client/proto"
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	num1 = flag.Int64("num1", 1, "num1的值")
	num2 = flag.Int64("num2", 2, "num2的值")
)

func main() {
	flag.Parse() // 从命令行解析num1，num2
	// 连接rpc server
	cc, err := grpc.Dial("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed:%v\n", err)
	}
	defer cc.Close()

	// 创建rpc客户端
	c := proto.NewCalcuClient(cc)
	// 调用add方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &proto.Input{Num1: (*num1), Num2: (*num2)})
	if err != nil {
		log.Fatalf("add failed:%v\n", err)
	}
	log.Printf("Answer is:%d", r.GetAns())

}
