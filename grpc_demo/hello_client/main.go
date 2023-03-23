package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hello_client/proto"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// grpc 客户端
// 调用server端的 SayHello 方法

var name = flag.String("name", "csy", "tell me who you are")

func main() {
	flag.Parse()
	// 连接server
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc dail failed, err:%v\n", err)
		return
	}
	defer conn.Close()

	// 创建客户端
	c := proto.NewGreeterClient(conn) // 生成的Go代码

	// 调用rpc方法 1.普通的rpc调用
	/*	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
			if err != nil {
				log.Fatalf("c.Sayhello failed, err:%v\n", err)
				return
			}
			// 拿到rpc响应
			log.Printf("resp:%v\v", resp.GetReply())*/

	// 流式rpc（服务端返回流）
	//callLotsOfReplies(c)

	// 流式rpc（客户端发送流）
	//callLotsOfGreetings(c)

	// 流式rpc（双向流）
	runBidiHello(c)
}

func callLotsOfReplies(c proto.GreeterClient) {
	// server端流式RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 2.流式rpc调用
	stream, err := c.LotsOfReplies(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err:%v\n", err)
	}
	for {
		// 接收服务端返回的流式数据，当收到io.EOF或错误时退出
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("recv failed, err:%v\n", err)
		}
		log.Printf("get replay:%q\n", recv.GetReply())
	}
}

func callLotsOfGreetings(c proto.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 客户端要流式的发送消息
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Printf("c.LotsOfGreetings(ctx) failed, err:%v\n", err)
	}
	names := []string{"lily", "lucy", "jenny"}
	for _, name := range names {
		// 流式的发送数据
		err := stream.Send(&proto.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v\"", name, err)
		}
	}
	// 流式发送之后要关闭流，然后接收返回值并打印
	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("c.LotsOfGreetings failed: %v", err)
	}
	log.Printf("get replay:%v", recv.GetReply())
}

func runBidiHello(c proto.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	// 双向流模式
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalf("c.BidiHello failed, err: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			// 接收服务端返回的响应
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("c.BidiHello stream.Recv() failed, err: %v", err)
			}
			fmt.Printf("AI：%s\n", in.GetReply())
		}
	}()
	// 从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	for {
		cmd, _ := reader.ReadString('\n') // 读到换行
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		// 将获取到的数据发送至服务端
		if err := stream.Send(&proto.HelloRequest{Name: cmd}); err != nil {
			log.Fatalf("c.BidiHello stream.Send(%v) failed: %v", cmd, err)
		}
	}
	stream.CloseSend()
	<-waitc
}
