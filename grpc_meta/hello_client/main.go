package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"code.xxx.com/backend/hello_client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "Siyu", "通过-name告诉server你是谁")

func main() {
	flag.Parse() // 解析命令行参数

	// 连接server
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
		return
	}
	defer conn.Close()
	// 创建客户端
	c := proto.NewGreeterClient(conn) // 使用生成的Go代码
	// 调用RPC方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 发起普通rpc调用
	// 带元数据
	md := metadata.Pairs(
		"token", "app_test_csy",
	)
	// 将元数据发送到服务端, NewOutgoingContext 将元数据附加到context
	ctx = metadata.NewOutgoingContext(ctx, md)
	// 在发起rpc调用前，声明两个变量
	var header, trailer metadata.MD
	resp, err := c.SayHello(
		ctx,
		&proto.HelloRequest{Name: *name},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		log.Printf("c.SayHello failed, err:%v", err)
		return
	}
	// 拿到响应数据之前，获取header
	fmt.Printf("header:%v\n", header)
	// 拿到了RPC响应
	log.Printf("resp:%v\n", resp.GetReply())
	// 拿到响应数据之后，获取trailer
	fmt.Printf("trailer:%v\n", trailer)

}
