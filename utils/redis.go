package utils

import (
	"Graduation/global"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisConnect() *redis.Client {
	var r = global.GVA_CONFIG.Redis
	options := redis.Options{
		Addr:            r.RedisHost + ":" + r.RedisPort,
		DB:              r.RedisDb,
		PoolSize:        r.RedisPoolSize,                                 // Redis连接池大小
		MaxRetries:      r.RedisMaxRetries,                               // 最大重试次数
		ConnMaxIdleTime: time.Second * time.Duration(r.RedisIdleTimeout), // 空闲链接超时时间
	}
	if r.RedisPassword != "" {
		options.Password = r.RedisPassword
	}
	Rdb := redis.NewClient(&options)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := Rdb.Ping(ctx).Result()
	if err == redis.Nil {
		global.GVA_LOG.Debug("[StoreRedis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("[StoreRedis] redis connRdb err,err=%s", err))
		panic(err)
	} else {
		global.GVA_LOG.Debug(fmt.Sprintf("[StoreRedis] redis connRdb success,suc=%s", pong))
	}
	return Rdb

}
