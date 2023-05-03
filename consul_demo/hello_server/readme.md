## consul服务注册
1. 启动创建一个grpc服务
2. 注册consul健康检查 （定时向节点发送请求，返回ok，类似ping-pong）
3. 连接至consul `api.NewClient(api.DefaultConfig())`
4. 获取本机出口ip （consul在容器内，并不是127.0.0.1）
5. 把服务注册到consul上
   1. 配置健康检查 （定义服务类型、超时时间、探测间隔、注销信息等）
   2. 定义一个consul服务体 `api.AgentServiceRegistration`
   3. 注册 `client.Agent().ServiceRegister(srv)`

## consul服务注销
1. 创建一个退出通道`quitCh`，阻塞监测退出信号
2. `cc.Agent().ServiceDeregister(serviceID)`
3. 启动服务要放在goroutine中