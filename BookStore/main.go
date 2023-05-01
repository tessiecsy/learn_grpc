package main

import (
	"BookStore/pb"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	// 连接数据库
	db, err := NewDB("bookstore.db")
	if err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
		return
	}

	// 创建server
	srv := server{
		bs: &bookstore{db: db},
	}
	// 启动grpc服务
	l, err := net.Listen("tcp", ":8091")
	if err != nil {
		fmt.Printf("failed to listen, err:%v\n", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterBookstoreServer(s, &srv) // 注册服务

	// 同一端口分别处理grpc与http服务
	// 1.创建grpc-gateway mux
	gwmux := runtime.NewServeMux()
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterBookstoreHandlerFromEndpoint(
		context.Background(), gwmux, "127.0.0.1:8091", dops); err != nil {
		log.Fatalf("RegisterBookstoreHandlerFromEndpoint failed, err:%v\n", err)
		return
	}

	// 2.新建http mux
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// 3.定义http server
	gwServer := &http.Server{
		Addr:    "127.0.0.1:8091",
		Handler: grpcHandlerFunc(s, mux),
	}
	// 4.启动
	fmt.Println("serving on 127.0.0.1:8091...")
	gwServer.Serve(l)
}

// grpcHandlerFunc 将gRPC请求和HTTP请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
