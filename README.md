

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


## 依赖
go get  gorm.io/gorm@v1.30.0 <br>
go get  gorm.io/driver/mysql@v1.6.0 <br>
go get  github.com/gin-gonic/gin@v1.12.0 <br>
go get  github.com/gin-contrib/cors <br>
//go get  go get github.com/rs/cors@v1.11.1 <br>
go get  github.com/golang-jwt/jwt/v5@v5.3.0 <br>
go get  github.com/redis/go-redis/v9 <br>
go get github.com/spf13/viper@v1.12.0 <br>
go get go.uber.org/zap@v1.21.0 <br>
gopkg.in/natefinch/lumberjack.v2@v2.0 <br>

## 说明
数据库 mysql <br>
ORM gorm <br>
路由 gin <br>
跨域 cors  <br>
JWT go-jwt <br>
Redis go-redis <br>
配置管理 viper <br>
日志管理 zap <br>
日志轮换 lumberjack <br>
文档管理 swag 待添加<br> 

