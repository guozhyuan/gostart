package test_test

import (
	"context"
	"gostart/internal/service"
	"testing"
	"time"
)

//String , Hash , Set, List

func TestConnectRedis(t *testing.T) {
	defer service.RedisDB.Close()
	service.ReadConfig()
	service.ConnectRedis()
	ctx := context.Background()

	cmd := service.RedisDB.Ping(ctx)
	if _, err := cmd.Result(); err != nil {
		t.Fatalf("Redis ping 失败: %v", err)
	}
	t.Log("Redis 连接成功")

	HSet(ctx)
	LPush(ctx)
}

func Set(c context.Context, key string, value string, expiration time.Duration) {
	service.RedisDB.Set(c, key, value, expiration)
}

func HSet(c context.Context) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	service.RedisDB.HSet(c, "HSet", m)
}

func LPush(c context.Context) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	service.RedisDB.LPush(c, "LPush", m)
	service.RedisDB.LPush(c, "LPush", m)
	service.RedisDB.LPush(c, "LPush", m)
}
