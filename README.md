# myserver

A Go web service project skeleton.

## Directory layout

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

## Run

```bash
go run ./cmd/server
```

Server listens on `:8080` (override with `APP_ADDR`).

## Health check

```bash
curl http://localhost:8080/api/health
```

## Test

```bash
go test ./...
```

## 依赖
go get  gorm.io/gorm@v1.30.0
go get  gorm.io/driver/mysql@v1.6.0
go get  github.com/gin-gonic/gin@v1.12.0

go get  github.com/gin-contrib/cors
//go get  go get github.com/rs/cors@v1.11.1

go get  github.com/golang-jwt/jwt/v5@v5.3.0
go get  github.com/redis/go-redis/v9

go get github.com/spf13/viper@v1.12.0

go get go.uber.org/zap@v1.21.0
gopkg.in/natefinch/lumberjack.v2@v2.0

// swag