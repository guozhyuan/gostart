

## 目录结构 


```
cmd/server      入口，main 函数
internal        私有代码，禁止外部引用
  ├── config    配置加载
  ├── handler   HTTP handler
  ├── middleware 中间件
  ├── model     数据模型
  ├── router    路由注册
  └── service   业务逻辑
pkg/response    可复用响应封装
api/docs        API 文档 (OpenAPI 等)
configs         配置文件
scripts         构建/运维脚本
test            集成测试
web             前端资源
```


## 说明
### ORM  mysql gorm 
```
go get  gorm.io/gorm@v1.30.0 
go get  gorm.io/driver/mysql@v1.6.0 
go get gorm.io/gen@v0.3.28    // 生成模型和CRUD
```
### 路由   gin 
```
go get  github.com/gin-gonic/gin@v1.12.0 
```
### 跨域   cors  
```
go get  github.com/gin-contrib/cors 
//go get  go get github.com/rs/cors@v1.11
```
### JWT   go-jwt 
```
go get github.com/golang-jwt/jwt/v5@v5.3.0 
```
### Redis   go-redis 
```
go get github.com/redis/go-redis/v9 
```
### 配置管理   viper 
```
go get github.com/spf13/viper@v1.12.0 
```
### 日志管理 zap & 日志轮换 lumberjack
```
go get go.uber.org/zap@v1.21.0 
go getgopkg.in/natefinch/lumberjack.v2@v2.0 
引入需要 import ""gopkg.in/natefinch/lumberjack.v2""
```
### 文档管理 swag 

```
<!-- go 要保证install和 get的包一致，否则会生成的docs.go可能会报错 --> 
go install github.com/swaggo/swag/cmd/swag@v1.16.1 
go get github.com/swaggo/swag@v1.16.1
go get github.com/swaggo/gin-swagger@v1.6.1 
go get github.com/swaggo/files
```

### 限流 
```
1.go.uber.org/ratelimit


2.go-redis/redis_rate


```

### 定时任务
```
go get github.com/robfig/cron/v3@v3.0.0
```



#  优秀go库
```
https://github.com/avelino/awesome-go
```
# Web脚手架
```
https://github.com/piupuer/gin-web
```

# 微服务框架
```
Kratos
go-zero
go-micro
```