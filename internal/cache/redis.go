// internal/cache/redis.go
package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"my-go-microservice/internal/config"
)

// NewRedisPool 根据配置创建 Redis 连接池
func NewRedisPool(cfg *config.Redis) *redis.Pool {
	return &redis.Pool{
		// 最大空闲连接数
		MaxIdle: cfg.Pool.MaxIdle,
		// 最大活跃连接数（0 表示无限制）
		MaxActive: cfg.Pool.MaxActive,
		// 空闲连接超时时间
		IdleTimeout: time.Duration(cfg.Pool.IdleTimeout) * time.Second,
		// 获取连接时是否等待
		Wait: cfg.Pool.Wait,

		// 连接工厂函数
		Dial: func() (redis.Conn, error) {
			// 设置连接超时
			c, err := redis.Dial(
				"tcp",
				cfg.Addr,
				redis.DialConnectTimeout(time.Duration(cfg.DialTimeout)*time.Second),
				redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Second),
			)
			if err != nil {
				return nil, err
			}

			// 认证密码（如果配置）
			if cfg.Password != "" {
				if _, err := c.Do("AUTH", cfg.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

			// 选择数据库
			if _, err := c.Do("SELECT", cfg.DB); err != nil {
				_ = c.Close()
				return nil, err
			}

			return c, nil
		},

		// 健康检查：PING 命令验证连接有效性
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// GetConnection 从连接池获取一个连接
func GetConnection(pool *redis.Pool) redis.Conn {
	return pool.Get()
}

// Ping 测试 Redis 连通性
func Ping(pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err
}
