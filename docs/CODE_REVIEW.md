# gostart 项目审查报告

审查范围：全部 Go 源码、配置、测试、Git 跟踪状态。基线为 `git status` 时的 HEAD (`6047b69`) 加 4 个未提交修改。

验证结果：

| 检查项 | 结果 |
| --- | --- |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过，无告警 |
| `go test ./...` | **失败**：`internal/pkg` 断言失败；`test` 包 panic |

结论：**能编译、但当前状态下不能通过测试，且存在若干可被直接利用的安全缺陷。**

> 审查期间工作区发生了并发改动：`internal/config/config.go` 与 `configs/config.yaml` 的 JWT 键名被先后修改（曾出现 YAML 用 `secretKey`、Go 结构体标签用 `secret` 的**互相矛盾**状态，会导致 viper 映射出空密钥）。问题 1 已按 `secretKey` 对齐两侧并完成修复，详见下文；其余小节依据审查时的实际内容编写，若你在此期间改过对应文件，请以代码为准。

---

## 一、严重（安全，建议立即处理）

### 1. JWT 密钥硬编码，配置项完全失效 —— ✅ 已修复

**原问题**：`internal/pkg/token.go` 把签名密钥写死为字面量 `"secret"`：

```go
token, err := jwtToken.SignedString([]byte("secret"))
...
jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) { return []byte("secret"), nil })
```

任何拿到源码的人都能自行签发任意 `userId` 的合法 token，认证体系形同虚设。同时配置链路也是断的：`internal/config/config.go` 的结构体标签与 `configs/config.yaml` 的键名不一致，viper 会映射出**空密钥**。

**已完成的修复**（`internal/pkg/token.go`、`internal/config/config.go`、`cmd/server/main.go`）：

- 密钥改为运行时读取 `config.Configs.JWT.SecretKey`，两侧键名统一为 `secretKey`（YAML 与结构体标签一致）。
- **不设任何兜底默认密钥**：配置缺失时 `GenerateToken`/`ParseAndValidateToken` 返回 `ErrJWTSecretNotConfigured`，而不是回退到内置常量。
- `cmd/server/main.go` 启动期 fail-fast：密钥为空直接 `log.Fatal` 拒绝启动，避免带病运行到首次登录才暴露。
- 顺带修掉同一函数里的两处关联缺陷：签名算法锁定为 HMAC（防 `alg=none`/算法混淆绕过）；移除 `fmt.Println("AuthMiddleware token: ", token)` —— 把完整 token 打进日志等同于泄露可用凭据。

**验证**：`go build ./...` 与 `go vet ./...` 均通过；临时加入 6 个 token 单测（配置密钥往返、字面量 `"secret"` 下旧 token 必须失效、密钥轮换后旧 token 失效、空密钥报错、config 为 nil 不 panic、token 类型不匹配拒绝）全部 PASS，验证后已删除这些临时测试文件。另用一次性程序实测真实配置映射结果：`secretKey="your_jwt_secret" expire=7200`，映射成功。

**仍需你完成的部分**：`configs/config.yaml` 里的值仍是占位串 `"your_jwt_secret"`，**上线前必须换成高强度随机密钥**（如 `openssl rand -base64 48`），并从环境变量注入。另外该文件已被 Git 跟踪（见问题 2），密钥入库等同于公开。

### 2. 明文凭据随 Git 分发
`configs/config.yaml` **已被 git 跟踪**（`git ls-files` 可见），其中含 MySQL 口令 `jugg`、Redis 口令 `guo`、内网地址 `172.31.98.14`。

影响：口令已经进入仓库历史，改文件不等于撤销泄露，必须轮换凭据。
修复：轮换 MySQL/Redis 口令；把该文件改为 `config.example.yaml` 入库、真实配置由部署环境注入或本地忽略；用 `viper.AutomaticEnv()` + 环境变量覆盖。

### 3. 认证中间件存在多处失效点
`internal/middleware/auth.go`：

- **第 27–32 行的 `Bearer` 前缀校验被注释掉**，第 35 行直接把整个 `Authorization` 头当 token 解析。标准客户端发 `Authorization: Bearer xxx` 会 401；而"能通过"的调用方式是非标准的裸 token。这是必然踩的兼容性坑，应恢复前缀剥离（同时容忍裸 token 只会放大歧义）。
- **白名单按完整路径精确匹配**（`slices.Contains(..., path)`，第 16 行）。`configs/config.yaml:38-41` 里的 `/api/streamer` 在路由表中是 `POST /api/streamer`，虽然此例恰好能匹配上，但模式很脆弱：任何带路径参数的接口（如 `/api/user/:id`）永远无法加入白名单。建议改用路由 pattern（`c.FullPath()`）匹配。
- **忽略签发者与受众**：`ParseAndValidateToken` 未校验 `iss`/`aud`，也没有任何吊销机制，而 `token.go:19` 的有效期是 30 天。登出接口 `internal/handler/user_handler.go:78` 是空实现，无法使已泄露 token 失效。建议引入 `jti` + Redis 黑名单（项目已有 `RedisDB`），并把有效期与配置里的 `jwt.expire` 对齐（当前配置 `7200` 与代码里写死的 `30 * 24 * time.Hour` 不一致，配置项同样形同虚设）。
- **`c.Set("userId", ...)` 后无人校验归属**：见下条。

### 4. 越权读取：可枚举任意用户
`internal/handler/user_handler.go:159` 取出 token 中的 `userId`，仅用 `fmt.Println` 打印，随后（第 161 行）按 URL 中的 `:id` 查询并返回，从不比对两者。

影响：任意登录用户可遍历 `/api/user/1..N` 拿到他人资料（含邮箱）。建议：普通用户接口应固定使用 token 中的 `userId`，或显式做 `claims.UserId == id` 判定；管理类查询单独放管理员路由并加角色校验。

### 5. 管理员登录是桩实现
`internal/handler/admin_handler.go:26-29` 打印一行日志后直接返回 `{"message": "Admin login successful"}`，无任何凭据校验。`/admin/login` 还是 `GET`。结合 `internal/router/router.go:35` 的限流（每 IP 令牌桶为全局共享，见问题 12），这是一个对外暴露的"永远成功"的登录端点。上线前必须补齐或删除。

### 6. bcrypt 哈希可能被序列化输出
`internal/dao/model/user.gen.go:17` 的 `Password` 字段标签为 `json:"password"`，没有 `json:"-"`。目前 handler 用手工组装的 `UserResp` 规避了泄露，但这属于"靠约定不出错"：一旦有人直接返回 `*model.User`，口令哈希即随响应泄露。建议在该字段加 `json:"-"`。

---

## 二、高（正确性与健壮性）

### 7. 启动链路 panic 风险（测试已实证）
`internal/config/config.go:92-98`：读配置失败只 `fmt.Println`，不返回错误，`Configs` 仍被赋值；而 `ReadConfigByViper()` 内部是 `viper.Unmarshal(Configs)`，字段全为零值、`Gin`/`Mysql`/`Zap` 等**指针字段为 nil**。

后果链：配置缺失 → `internal/pkg/init.go:31` 访问 `config.Configs.Mysql.User` 时 **nil 解引用 panic**。`go test ./...` 中的 `Test_insertPlatform` 正是这样崩的（未先读配置就调用 `pkg.ConnectDB()`，栈显示 `init.go:31`）。

同类问题：`internal/pkg/init.go:55-58` 的 `ConnectRedis()` 只构造客户端、不 `Ping`，Redis 不可用时故障被推迟到第一个业务命令。

修复：配置加载失败应返回 error 并由 `main` 终止（fail-fast）；把 `Config` 结构体的指针字段改为值类型，或在初始化时显式给默认值；`ConnectRedis` 增加 `Ping` 探活；`ConnectDB` 里 `sqlDB, _ := db.DB()` 的错误也被丢弃，应一并处理。

### 8. `main` 忽略监听错误、无优雅退出
`cmd/server/main.go:32` `engine.SetTrustedProxies(...)` 的返回值被丢弃；`main.go:41` `s.ListenAndServe()` 的返回值同样被丢弃。

后果：端口被占用时进程静默退出且退出码为 0，容器/编排系统会误判为正常退出而不重启，排障时看不到任何线索。

`main.go:44-48` 的 `Release()` 是空壳且**从未被调用**。`internal/pkg/init.go:137` 的 `defer logger.Sync()` 在 `ZapLogInit` 返回时就已执行，起不到进程退出时刷盘的作用 —— 崩溃或强杀时最后一批日志会丢失。

修复：`ListenAndServe` 的错误写入日志并以非零码退出；注册 `SIGINT/SIGTERM` 做 `s.Shutdown(ctx)` 宽限关闭，并在关闭流程里执行 `sqlDB.Close()`、`RedisDB.Close()`、`logger.Sync()`。项目笔记 `notes/常见问题.md:132-137` 已经记录了这个模式，落地即可。

### 9. 日志中间件注册顺序错误，401 不落盘
`internal/router/router.go:19-20`：

```go
api.Use(middleware.AuthMiddleware())
api.Use(middleware.ZapLoggerMiddleware())
```

Gin 中先声明的中间件先执行；`AuthMiddleware` 认证失败时调用 `c.Abort()`，其后的 `ZapLoggerMiddleware` 直接跳过。**所有认证失败请求都不会产生日志。** 而 401 恰恰是安全审计最需要的事件。

修复：把 `ZapLoggerMiddleware` 放到 `AuthMiddleware` 之前。

### 10. 缺 `return` / 忽略错误的响应路径
- `internal/handler/user_handler.go:155-157`：`strconv.Atoi` 失败后写了 400 但没有 `return`，继续执行到第 172 行又写一次响应（Gin 会记录 `headers were already written`），并带 `intId == 0` 查库。
- `internal/handler/user_handler.go:229-232`：`s, oerr := strconv.Atoi(size)` 出错后同样缺 `return`；更糟的是这里复用了变量 `oerr`，把前一个 `offset` 的错误覆盖掉了，`if oerr != nil` 的语义已经被混淆。
- `internal/handler/user_handler.go:127-139`：`GetUsers` 先遍历 `ret` 组装响应，再在**之后**判断 `err != nil`；出错时 `ret` 为 nil，遍历无害但顺序应反过来。
- `internal/service/user_service.go:100`：查 `platform` 时丢弃错误（`platform, _ := ...`），查询失败会静默返回 `Platform: nil`。
- `internal/service/user_service.go:78-80`：`user.ID == 0` 的判断永远不可达 —— 用 gen 的 `First()` 查不到记录时返回的是 `gorm.ErrRecordNotFound` 错误，已在上一行返回。

### 11. N+1 查询
`internal/service/user_service.go:99-108`：循环内逐条 `query.Platform.Where(...).First()`，一页 N 条主播就产生 N+1 次往返，且失败被吞（见上条）。应改为一次 `IN` 查询后建 `map[platformID]*Platform` 关联，或使用 gen 的关联/`Preload` 能力。

### 12. 限流语义不当且只保护了一个桩接口
`internal/middleware/ratelimiter.go:23` 用 `ratelimit.New(rate)` 创建**单个全局令牌桶**，`Take()` 在超速时**阻塞排队**而不是快速拒绝。

后果：并发高峰时请求会堆积在中间件里等待（放大尾延迟、占用连接），且限流额度被所有客户端共享，单个恶意客户端即可耗尽全站额度、拖垮正常用户。理想做法是按 `ClientIP` 或用户维度建桶，超限立即返回 `429`。另外限流只挂在 `/admin` 组（且该组只有桩接口），真正需要防爆破的 `POST /api/login`、`POST /api/regist` 反而没有任何限流。

### 13. CORS 配置细节
`internal/middleware/cors.go:14` 的 `AllowOrigins` 里 `"http://172.31.98.14"` 重复出现两次（显然是从 `"https://..."` 复制后漏改）；`AllowCredentials: true` 配白名单本身没问题，但白名单是硬编码在代码里的，无法按环境区分 dev/prod。建议移入配置。`AllowMethods` 未包含 `PATCH`/`HEAD`。

---

## 三、中（工程化与一致性）

### 14. 测试几乎不可运行、且当前是红的
- `internal/pkg/crypt_test.go:12` 用一个**伪造的哈希常量**（`$10$0BX4g0YCr9tANo8i4n2naOa2FyPPnPEiN3GqxuZGTtjUketXyFhoe`）去 `Compare("guo")`，必然失败。它想验证的是常量，不是 `Encrypt`/`Compare` 的行为 —— 应改为 `Compare(hash, "guo") == true` 且 `Compare(hash, "wrong") == false`。
- `test/redis_test.go:13-15`：`defer pkg.RedisDB.Close()` 在 `pkg.ConnectRedis()` **之前**注册，此时 `RedisDB` 还是 nil，一旦 panic 就是二次 panic；且该测试硬依赖 `172.31.98.14:6379` 可达。
- `test/migrate.go:30-38` 里 `ReadConfigByViper`+`ConnectDB` 被连续调用了两遍。
- `test/migrate.go:84-116` 的 `insertPlatformStreamer` 请求外部站点 `https://live.suancaihu.eu.org/...` 并向数据库写入，且 `http.Get` 出错后只打印不 return，紧接着 `resp.Body` 会 nil 解引用。
- `internal/middleware/auth_test.go` 只有一行 `package middleware`，是空文件。
- **核心逻辑零覆盖**：认证中间件、`GenerateToken`/`ParseAndValidateToken`、`service` 层都没有测试。

建议：把需要外部依赖的验证移出 `go test`（放到 `scripts/` 或构建 tag 后面），为认证与 token 补纯内存单测，并在 CI 里把 `go test ./...` 变成必须绿的关卡。

### 15. 无效/冗余代码应当清理或接上
- `internal/handler/user_handler.go:38-40`：`Test` 空实现；`:78-80` `Logout` 空实现；`:176-203` `UpdateUser`/`DeleteUser` 整段被注释，但接口（`:22-23`）和路由（`router.go:28-29`）仍然暴露 `PUT`/`DELETE /api/user/:id`，调用会得到 200 空响应体 —— 这是比 404 更容易误导调用方的行为。
- `internal/config/config.go:68-82` `ReadConfig2()` 与 `:84-100` `ReadConfigByViper()` 并存，只有后者被用（`internal/pkg/init.go:20`）。
- `internal/middleware/log.go:17` `Logging` 已标 `Deprecated`、`:81` `LogMiddleware` 未被引用；`:101-104` `outputWriter()` 每次调用 `os.OpenFile` 后**从不关闭**，是文件描述符泄漏。若要保留，改成包级一次性初始化并在退出时关闭。
- `internal/dao/common/user_resp.go:5-12`：`LoginResp` 是 HTTP 响应体，却带着 `gorm:"primarykey"` 等持久化标签，且 `AccessToken`/`RefreshToken` 都标了 `primarykey` —— 这些标签在此处毫无意义且误导。
- 与之对照，`internal/dao/model/*.gen.go` 顶部标注 "Code generated by gorm.io/gen. DO NOT EDIT."，这些文件**不应手工修改**；自定义响应类型放在 `dao/common` 是合理的，但目录命名（`common` 同时放 model 与 resp）会随规模增长变得含糊。

### 16. README 与实际结构不符
`README.md:6-21` 描述的目录树**在仓库中不存在**：

- 它写 `internal/model`、`pkg/response`、`api/docs`；实际是 `internal/dao/model`、`internal/dao/common`、`docs`。
- `internal/service`、`internal/config` 存在（这部分正确），但没有提到 `internal/pkg`、`internal/dao/query`（gen 生成的 CRUD，是数据访问的真正入口）。
- `web/`、`scripts/` 目前只有 `.gitkeep` 占位。

建议按真实结构重写这段，并补一节"如何跑起来"（依赖 MySQL/Redis、`configs/config.yaml` 必填字段、`go run ./cmd/server`、Swagger 地址 `http://localhost:9090/swagger/index.html`）。README 里 `@host localhost:8080` 与配置端口 `9090`（`configs/config.yaml:5`）不一致，`cmd/server/main.go:16` 的 swag 注释同样写 `localhost:8080` —— 文档与实际会互相误导。

### 17. 其他小项
- `cmd/server/main.go:32`：`TrustProxy` 在生产必须收敛到真实反代地址，否则 `ClientIP()` 可被 `X-Forwarded-For` 伪造，连带影响限流与日志（项目已把 `172.31.98.14` 列入白名单，思路是对的，注意 prod 环境别放行整个网段）。
- DSN 拼接（`internal/pkg/init.go:30-37`）没有转义口令中的 `@`、`:`、`/` 等字符，口令一旦包含这些字符连接串会被解析错位；另 `Collation` 配置项定义了却未使用。
- `cmd/migrate/migrate.go:44-48` 与 `test/migrate.go:118-124` 的 `migrate()` 均为空，自动建表能力实际不可用。若依赖 gorm auto-migrate，应接上；否则建议删除以免误解。
- `test/migrate.go:13` 的 `Platform.ID` 用 `uint64`，`internal/dao/model/platform.gen.go:11` 生成的是 `int32`，同一张表两套类型定义，容易在后续代码里埋类型不一致的坑。
- `go.mod:1` 模块名为裸 `gostart`，无域名前缀，无法作为可导入模块发布。
- `go.mod:16-17` 的 `zap v1.21.0`、`swag v1.16.1` 相对陈旧（zap 已到 v1.27+），非阻塞项，可在方便时升级。
- 仓库根目录有 58 MB 的 `main.exe` 构建产物（已被 `.gitignore` 的 `*.exe` 忽略，未入库，这点做得对）；`.git` 目录 27 MiB，说明历史上曾提交过二进制，清理 `git gc` 之前需谨慎处理历史。

---

## 四、修复优先级建议

**第一梯队（安全，当天）**
1. ~~JWT 密钥配置化~~ ✅ 已完成（代码层）；**待你执行**：轮换泄露的 MySQL/Redis 凭据，并把 `your_jwt_secret` 替换为强随机密钥（问题 1、2）
2. 恢复 `Bearer` 前缀处理；`GetUser` 增加归属校验；删除或补齐 `/admin/login` 桩（问题 3、4、5）
3. `Password` 字段加 `json:"-"`（问题 6）

**第二梯队（稳定性，本周）**
4. 配置加载 fail-fast、`ConnectDB/ConnectRedis` 错误处理与探活（问题 7）
5. `ListenAndServe` 错误处理 + 优雅退出（问题 8）
6. 调整日志中间件顺序；补齐缺失 `return`；修掉 N+1（问题 9、10、11）

**第三梯队（工程化）**
7. 修复红色测试、为认证/token/service 补单测、把外部依赖测试移出 `go test`（问题 14）
8. 清理空实现与死代码、限流改为按 IP 拒绝式并覆盖登录接口、CORS/端口/README 对齐（问题 12、13、15、16）

---

## 五、做得好的地方

审查不应只有批评，以下几点值得保持：

- 分层清晰：`handler`（薄）→ `service`（业务编排）→ `dao/query`（gen 生成的 CRUD），`internal/service/user_service.go:14-16` 还用注释把职责写清楚了。
- 使用 gorm/gen 生成类型安全的查询代码，避免手写字符串条件，`internal/dao/query/user.gen.go` 即为其产物。
- 口令使用 bcrypt（`internal/pkg/crypt.go`），没有自己造哈希。
- 日志按环境分流（dev 带彩色控制台 + 文件，prod 纯文件），并用 lumberjack 做轮转（`internal/pkg/init.go:82-96`），这套配置比多数脚手架完整。
- `go vet` 干净，说明基础代码卫生不错。
- 有意识地记录笔记（`notes/`）与使用 Swagger 注释，`internal/handler/user_handler.go` 的注解格式规范。
- 构建产物与日志目录已正确写入 `.gitignore`（`main.exe`、`logs/` 确实未被跟踪）。

---

*报告基于静态审查 + `go build`/`go vet`/`go test` 实际执行结果；未连接真实的 MySQL/Redis 实例做运行时验证。*
