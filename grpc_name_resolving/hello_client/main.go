package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"hello_client/pb"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "q1mi", "通过-name告诉server你是谁")

func main() {
	flag.Parse() // 解析命令行参数

	// 指定连接server
	//conn, err := grpc.Dial("127.0.0.1:8972",
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)

	// dns解析
	//conn, err := grpc.Dial("dns:///127.0.0.1:8972",
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)

	// 自定义解析
	//grpc.Dial函数中会先根据q1mi这个scheme找到我们通过init函数注册的q1miResolverBuilder，
	//然后调用它的Build()方法构建我们自定义的q1miResolver，并调用ResolveNow()方法获取到服务端地址。
	conn, err := grpc.Dial("q1mi:///resolver.liwenzhou.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithResolvers(&q1miResolverBuilder{}),  // 此时不用写init函数
	)

	// 自定义解析 + 轮询的负载均衡
	//conn, err := grpc.Dial("q1mi:///resolver.liwenzhou.com",
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 这里设置初始策略
	//	//grpc.WithResolvers(&q1miResolverBuilder{}),
	//)
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
		return
	}

	defer conn.Close()
	// 创建客户端
	c := pb.NewGreeterClient(conn)
	// 调用RPC方法
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			fmt.Printf("c.SayHello failed, err:%v\n", err)
			return
		}
		// 拿到了RPC响应
		fmt.Printf("resp:%v\n", resp.GetReply())
	}
}
