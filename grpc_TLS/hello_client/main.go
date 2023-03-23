package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"time"

	"code.xxx.com/backend/hello_client/proto"
	"google.golang.org/grpc"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "Siyu", "通过-name告诉server你是谁")

func main() {
	flag.Parse() // 解析命令行参数

	// 连接server
	// 加载证书
	cerds, err := credentials.NewClientTLSFromFile("cert/server.crt", "liwenzhou.com")
	if err != nil {
		fmt.Printf("credentials.NewClientTLSFromFile failed, er:%v\n, err")
		return
	}
	// 把证书传进去
	conn, err := grpc.Dial("127.0.0.1:8972",
		grpc.WithTransportCredentials(cerds),
	)
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
	resp, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Printf("c.SayHello failed, err:%v", err)
		return
	}
	// 拿到了RPC响应
	log.Printf("resp:%v", resp.GetReply())

}
