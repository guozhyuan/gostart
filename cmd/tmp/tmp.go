package main

import (
	"myserver/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// ret, _ := pkg.Encrypt("guo")
	// fmt.Println(ret)
	// fmt.Println(pkg.Compare(ret, "guo")) // true
	config.ReadConfigByViper()

	ZapLogInit()
	Zap.Debug("hello zap!")
	Zap.Info("hello zap!")
	Zap.Warn("hello zap!")
	Zap.Error("hello zap!")
	Zap.Fatal("hello zap!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	/* service.ZapInit()

	bytes, _ := json.Marshal(config.Configs)
	fmt.Println(string(bytes))

	service.ZapLogger.Info("配置信息", zap.Any("Zaaaap", config.Configs)) */

	/*
		 service.ConnectRedis()

		redis.SetLogLevel(3)
		var statusCmd *redis.StatusCmd = service.RedisDB.Set(service.Ctx, "username:guo", "guo content", 1*time.Hour)
		if statusCmd.Err() != nil {
			fmt.Println(statusCmd.Err())
		} else {
			fmt.Println(statusCmd)
		}


		// resp := LoginResp{
		// 	ID:           10,
		// 	Username:     "guo",
		// 	Email:        "guo@163.com",
		// 	Age:          20,
		// 	AccessToken:  "token",
		// 	RefreshToken: "refreshToken",
		// 	Address: &Address{
		// 		Addr: "addr_",
		// 		City: "city_",
		// 	},
		// }
		// bs, _ := json.Marshal(resp)
		service.RedisDB.Set(service.Ctx, "user:id:10", "first add", 30*7*24*time.Hour)
		//service.RedisDB.SetNX(service.Ctx, "user:id:10", "second add", 30*7*24*time.Hour)
		service.RedisDB.Set(service.Ctx, "user:id:10", "second add", 30*7*24*time.Hour)

		// Hash : HSet手动设置过期时间
		service.RedisDB.HSet(service.Ctx, "map:user:10", map[string]string{"key1": "hello", "key2": "world", "key3": "end"})
		service.RedisDB.Expire(service.Ctx, "map:user:10", 30*7*24*time.Hour)

		// List  LPop
		service.RedisDB.RPush(service.Ctx, "list:user:10", "item1", "item2", "item3")
		service.RedisDB.Expire(service.Ctx, "list:user:10", 30*7*24*time.Hour)

		// Set  SIsMember:检查是否存在  SMembers获取所有
		service.RedisDB.SAdd(service.Ctx, "set:user:10", "python", "redis", "database", "database")

		// SortedSet

		// 管道
		pipe := service.RedisDB.Pipeline()
		pipe.HSet(service.Ctx, "map:user:10", map[string]string{"key1": "hello", "key2": "world", "key3": "end"})
		pipe.SAdd(service.Ctx, "set:user:10", "python", "redis", "database", "database")
		pipe.RPush(service.Ctx, "list:user:10", "item1", "item2", "item3")

		service.RedisDB.Close()

	*/
}

type LoginResp struct {
	ID           uint64   `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Age          uint8    `json:"age"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Address      *Address `json:"address"`
}
type Address struct {
	Addr string `json:"addr"`
	City string `json:"city"`
}
