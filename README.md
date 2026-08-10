

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
```
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
<!-- go 要保证install和 get的包一致，否则会生成的docs.go可能会报错 --> 
go install github.com/swaggo/swag/cmd/swag@v1.16.1 <br>
go get github.com/swaggo/swag@v1.16.1
go get github.com/swaggo/gin-swagger@v1.6.1 <br>
go get github.com/swaggo/files <br>
```
## 说明
### 数据库   mysql <br>
### ORM   gorm <br>
### 路由   gin <br>
### 跨域   cors  <br>
### JWT   go-jwt <br>
### Redis   go-redis <br>
### 配置管理   viper <br>
### 日志管理   zap <br>
### 日志轮换 lumberjack <br>
### 文档管理 swag 

参考:https://github.com/swaggo/swag/blob/master/README_zh-CN.md
1.在main()添加以下注释, 
```
// @title           Gin + Swagger 示例 API
// @version         1.0
// @description     这是一个用 Gin 框架集成的 Swagger 示例
// @host            localhost:8848
// @BasePath        /
```
2.添加swagger路由
```
engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
```
3.生成文档
```
swag init -g ./cmd/server/main.go   <!-- 默认当前目录,用 -o 指定目录 -->
```
4.main.go导入生成的docs
```
import _ "gostart/docs" 
```
5.运行项目
```
go run cmd/server/main.go
```
6.访问文档
```
http://localhost:8848/swagger/index.html
```

#### 常见问题
授权验证
在main中添加
```
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
```
在具体handler中添加
```
// @Security ApiKeyAuth
```

