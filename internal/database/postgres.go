package database

import (
	"fmt"
	"my-go-microservice/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgreSQL(cfg *config.PostgresDB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// 获取底层 *sql.DB 对象以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大打开连接数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 连接最大空闲时间

	// 设置连接健康检查间隔
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

// 示例：带重试机制的数据库连接
func ConnectWithRetry(cfg *config.PostgresDB, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = InitPostgreSQL(cfg)
		if err == nil {
			return db, nil
		}

		// 指数退避：1s, 2s, 4s, 8s...
		backoff := time.Second * time.Duration(1<<i)
		time.Sleep(backoff)
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", maxRetries, err)
}
