package test_test

import (
	"context"
	"gostart/internal/pkg"
	"testing"
	"time"
)

//String , Hash , Set, List

func TestConnectRedis(t *testing.T) {
	defer pkg.RedisDB.Close()
	pkg.ReadConfig()
	pkg.ConnectRedis()
	ctx := context.Background()

	cmd := pkg.RedisDB.Ping(ctx)
	if _, err := cmd.Result(); err != nil {
		t.Fatalf("Redis ping 失败: %v", err)
	}
	t.Log("Redis 连接成功")

	HSet(ctx)
	LPush(ctx)
}

func Set(c context.Context, key string, value string, expiration time.Duration) {
	pkg.RedisDB.Set(c, key, value, expiration)
}

func HSet(c context.Context) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	pkg.RedisDB.HSet(c, "HSet", m)
}

func LPush(c context.Context) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	pkg.RedisDB.LPush(c, "LPush", m)
	pkg.RedisDB.LPush(c, "LPush", m)
	pkg.RedisDB.LPush(c, "LPush", m)
}
