package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	// 设置配置文件名称（不带扩展名）
	v.SetConfigName("config")
	// 支持的配置文件类型
	v.SetConfigType("yaml")
	// 配置文件搜索路径
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/myapp/")
	// 允许环境变量覆盖 - 使用更明确的环境变量前缀和分隔符
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 尝试读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 如果配置文件不存在，则继续（可能完全由环境变量驱动）
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		fmt.Printf("[DEBUG] Config file loaded from: %s\n", v.ConfigFileUsed())
	}

	// 将 Viper 配置绑定到结构体
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// 验证必要字段
	if err := validateConfig(c); err != nil {
		return nil, err
	}

	return &c, nil
}

func validateConfig(c Config) error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http.port must be a valid port number")
	}
	if c.GRPC.Port <= 0 || c.GRPC.Port > 65535 {
		return fmt.Errorf("grpc.port must be a valid port number")
	}
	if len(c.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka.brokers cannot be empty")
	}
	return nil
}
