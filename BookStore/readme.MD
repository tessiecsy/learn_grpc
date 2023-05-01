protoc -I=pb --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative bookstore.proto


## 添加分页功能

1. 把分页的代码添加进来
2. 修改 proto 文件
3. 重新生成代码

## grpc与 http使用同同一个端口
1. 创建grpc-gateway mux
2. RegisterBookstoreHandlerFromEndpoint
3. 新建http mux
4. 定义http server，处理函数grpcHandlerFunc处理两种不同请求
    