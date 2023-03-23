package main

import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hello_server/pb"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// grpc server

type server struct {
	pb.UnimplementedGreeterServer
	mu    sync.Mutex     // map并发锁
	count map[string]int // 存储每个name调用SayHello的次数
}

// SayHello 是我们需要实现的方法
// 这个方法是我们对外提供的服务
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count[in.GetName()]++ // 每次请求次数+1
	if s.count[in.GetName()] > 1 {
		// 返回请求次数限制的错误
		st := status.New(codes.ResourceExhausted, "request limit")
		// 添加错误的详情信息
		ds, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{
					{
						Subject:     fmt.Sprintf("name:%s", in.Name),
						Description: "每个name只能调用1次SayHello",
					},
				},
			},
		)
		if err != nil {
			// WithDetail执行失败了，返回原来的status.Err
			return nil, st.Err()
		}
		return nil, ds.Err() // 带Details的Err
	}
	// 正常执行
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
	s := grpc.NewServer() // 创建grpc服务
	// 注册服务
	pb.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		fmt.Printf("failed to serve,err:%v\n", err)
		return
	}
}
