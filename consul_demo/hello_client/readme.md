## consul服务发现
### case1
1. 连接consul
2. 根据服务名称查询服务实例 `cc.Agent().ServicesWithFilter("Service==hello")`
3. 返回服务实例map[string]*AgentService key为注册服务ID
4. 根据consul返回的服务实例建立连接

### case2
使用consul名称解析器  
`grpc.Dial("consul://xxxx:8500/hello",, grpc.WithTransportCredentials(insecure.NewCredentials()))`
