# 第一个gRPC示例

 hello world

## 三个步骤
1.编写protobuf文件
2.生成代码(服务端和客户端)
3.编写业务逻辑代码


## 作业
自己动手写一个 grpc版本的  add(x, y int) int


protoc -I=pb --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative hello.proto