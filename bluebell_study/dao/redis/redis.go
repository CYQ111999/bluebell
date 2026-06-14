package redis

import (
	"bluebell_study/setting"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init 初始化连接
func Init(cfg *setting.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}

// SetUserToken 将用户Token存入Redis，设置过期时间
func SetUserToken(userID int64, token string, expireTime time.Duration) error {
	key := fmt.Sprintf("user:%d", userID)
	return client.Set(key, token, expireTime).Err()
}

// GetUserToken 从Redis获取用户的Token
func GetUserToken(userID int64) (string, error) {
	key := fmt.Sprintf("user:%d", userID)
	return client.Get(key).Result()
}

// DeleteUserToken 从Redis删除用户的Token
func DeleteUserToken(userID int64) error {
	key := fmt.Sprintf("user:%d", userID)
	return client.Del(key).Err()
}
