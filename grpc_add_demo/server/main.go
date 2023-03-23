package main

import (
	"context"
	"log"
	"net"
	"server/pb"

	"google.golang.org/grpc"
)

// 兼容性，需要嵌入Unimplemented...
type server struct {
	pb.UnimplementedCalcuServer
}

// 实现service的Add方法
func (s *server) Add(cx context.Context, in *pb.Input) (*pb.Output, error) {
	return &pb.Output{Ans: in.GetNum1() + in.GetNum2()}, nil
}

func main() {
	// 在8972端口启动一个服务
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("failed to listen, err:%v\n", err)
		return
	}
	// 启动一个rpc服务
	s := grpc.NewServer()
	// 注册
	pb.RegisterCalcuServer(s, &server{})
	// 启动服务，把Listen的服务l返回
	err = s.Serve(l)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
		return
	}
}
