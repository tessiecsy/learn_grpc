package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"hello_server/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const serviceName = "hello"

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
	l, err := net.Listen("tcp", ":8976")
	if err != nil {
		fmt.Printf("failed to listen, err:%v\n", err)
		return
	}

	s := grpc.NewServer() // 创建grpc服务
	// 注册服务
	pb.RegisterGreeterServer(s, &server{})
	// 注册健康检查，给我们的grpc服务增加了健康检查的处理逻辑
	healthpb.RegisterHealthServer(s, health.NewServer()) // consul发来健康检查的要求，这个用户返回ok

	// 连接至consul
	config := api.DefaultConfig()
	config.Address = "192.168.6.128:8500"
	cc, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("api.NewClient failed, err:%v\n", err)
		return
	}
	// 获取本机出口ip
	ipinfo, err := GetOutboundIP()
	if err != nil {
		fmt.Printf("GetOutboundIP failed,err:%v\n", err)
		return
	}
	fmt.Printf("ip is:%s\n", ipinfo.String())
	// 把服务注册到consul上
	// 1.定义我们的服务
	// 配置健康检查，告诉consul如何进行健康检查
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ipinfo.String(), 8976), // 外网地址
		Timeout:                        "5s",
		Interval:                       "5s", // 间隔
		DeregisterCriticalServiceAfter: "1m", // 1minute后注销掉不健康的节点,文档默认最小1分钟，实际一般30s，时间尽量长一点，以防网络抖动
	}
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, "127.0.0.1", 8976)
	srv := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    []string{"tessie"},
		Port:    8976,
		Address: ipinfo.String(),
		Check:   check,
	}
	// 2.注册服务到consul
	cc.Agent().ServiceRegister(srv)

	// 启动服务
	// 要放到一个goroutine中，程序才能向下走到quitCh，否则是主线程一直在运行着处理client请求
	go func() {
		if err := s.Serve(l); err != nil {
			fmt.Printf("failed to serve,err:%v\n", err)
			return
		}
	}()
	// Ctrl+C 退出程序
	fmt.Println("wait to quit...")
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	<-quitCh
	// 程序退出的时候注销服务
	fmt.Println("service out...")
	err = cc.Agent().ServiceDeregister(serviceID)
	if err != nil {
		fmt.Printf("ServiceDeregister failed,err;%v\n", err)
		return
	}
}

// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80") // 尝试拨号拿到addr
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
