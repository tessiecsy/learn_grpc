package main

import (
	"code.xxx.com/backend/hello_client/proto"
	"context"
	"flag"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "七米", "通过-name告诉server你是谁")

func main() {
	flag.Parse() // 解析命令行参数

	/*
		// 1.连接consul
		config := api.DefaultConfig()
		config.Address = "192.168.6.128:8500"
		cc, err := api.NewClient(config)
		if err != nil {
			fmt.Printf("api.NewClient failed, err:%v\n", err)
			return
		}
		// 2.根据服务名称查询服务实例
		// map[string]*AgentService  key为服务ID
		serviceMap, err := cc.Agent().ServicesWithFilter("Service==`hello`") // 查询服务名称是hello的
		if err != nil {
			fmt.Printf("query `hello` service failed,err:%v\n", err)
			return
		}
		// 从consul返回的数据中选一个服务实例（机器）
		var addr string
		for k, v := range serviceMap {
			fmt.Printf("%s:%#v\n", k, v)
			addr = fmt.Sprintf("%s:%d", v.Address, v.Port) // 取第一个机器的addr和port来连接
			continue
		}
		log.Println(addr)
		// 3.根据consul返回的服务实例建立连接
		// 连接server
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	*/

	conn, err := grpc.Dial(
		"consul://192.168.6.128:8500/hello", // 使用consul名称解析器
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial failed,err:%v", err)
		return
	}
	defer conn.Close()
	// 创建客户端
	c := proto.NewGreeterClient(conn) // 使用生成的Go代码
	// 4.发起rpc调用
	// 调用RPC方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Printf("c.SayHello failed, err:%v", err)
		return
	}
	// 拿到了RPC响应
	log.Printf("resp:%v", resp.GetReply())

}
