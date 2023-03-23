package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials"
	"hello_server/pb"
	"net"
	
	"google.golang.org/grpc"
)

// grpc server

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello 是我们需要实现的方法
// 这个方法是我们对外提供的服务
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "hello " + in.GetName()
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	// 启动服务
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen, err:%v\n", err)
		return
	}
	// 加载证书信息
	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		fmt.Printf("credentials.NewServerTLSFromFile failed, err:%v\n", err)
		return
	}
	s := grpc.NewServer(grpc.Creds(creds)) // 创建grpc服务,让服务基于证书进行加密操作
	// 注册服务
	pb.RegisterGreeterServer(s, &server{})
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		fmt.Printf("failed to serve,err:%v\n", err)
		return
	}
}
