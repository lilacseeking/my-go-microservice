package config

import "time"

// Config represents the entire application configuration
type Config struct {
	App      App      `yaml:"app"`
	HTTP     HTTP     `yaml:"http"`
	GRPC     GRPC     `yaml:"grpc"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
	Kafka    Kafka    `yaml:"kafka"`
	Logger   Logger   `yaml:"logger"`
}

// App holds application-level settings
type App struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

// HTTP holds HTTP server settings
type HTTP struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// GRPC holds gRPC server settings
type GRPC struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	MaxRecvMsgSize int           `yaml:"max_recv_msg_size"`
	MaxSendMsgSize int           `yaml:"max_send_msg_size"`
	KeepaliveTime  time.Duration `yaml:"keepalive_time"`
	Timeout        time.Duration `yaml:"timeout"`
}

// Database holds database settings
type Database struct {
	Postgres PostgresDB `yaml:"postgres"`
}

// PostgresDB holds PostgreSQL connection details
type PostgresDB struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	Timezone        string        `yaml:"timezone"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

// Redis holds Redis client settings
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Pool     struct {
		MaxIdle     int           `yaml:"max_idle"`
		MaxActive   int           `yaml:"max_active"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
		Wait        bool          `yaml:"wait"`
	} `yaml:"pool"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
}

// Kafka holds Kafka client settings
type Kafka struct {
	Brokers  []string `yaml:"brokers"`
	Producer struct {
		RequiredAcks    int           `yaml:"required_acks"`
		Compression     string        `yaml:"compression"`
		MaxMessageBytes int           `yaml:"max_message_bytes"`
		Timeout         time.Duration `yaml:"timeout"`
	} `yaml:"producer"`
	Consumer struct {
		GroupID           string        `yaml:"group_id"`
		InitialOffset     string        `yaml:"initial_offset"`
		SessionTimeout    time.Duration `yaml:"session_timeout"`
		HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
		MaxProcessingTime time.Duration `yaml:"max_processing_time"`
	} `yaml:"consumer"`
	Topics struct {
		UserCreated    string `yaml:"user_created"`
		OrderProcessed string `yaml:"order_processed"`
	} `yaml:"topics"`
}

// Logger holds logging settings
type Logger struct {
	Level             string `yaml:"level"`
	Format            string `yaml:"format"`
	Output            string `yaml:"output"`
	EnableCaller      bool   `yaml:"enable_caller"`
	DisableStacktrace bool   `yaml:"disable_stacktrace"`
}
