package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"my-go-microservice/internal/config"
	"time"
)

// KafkaProducer 封装 Sarama 生产者客户端
type KafkaProducer struct {
	producer sarama.SyncProducer
}

// NewKafkaProducer 根据配置创建并启动 Kafka 生产者
func NewKafkaProducer(cfg *config.Kafka) (*KafkaProducer, error) {
	// 创建 Sarama 配置实例
	config := sarama.NewConfig()

	// 设置生产者模式为同步发送
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// 映射 YAML 配置到 Sarama 配置
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.Producer.RequiredAcks)
	config.Producer.Compression = getCompressionType(cfg.Producer.Compression)
	config.Producer.MaxMessageBytes = cfg.Producer.MaxMessageBytes
	config.Producer.Timeout = time.Duration(cfg.Producer.Timeout)

	// 创建同步生产者
	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

// getCompressionType 将字符串转换为 Sarama 压缩类型
func getCompressionType(compression string) sarama.CompressionCodec {
	switch compression {
	case "gzip":
		return sarama.CompressionGZIP
	case "snappy":
		return sarama.CompressionSnappy
	case "lz4":
		return sarama.CompressionLZ4
	case "zstd":
		return sarama.CompressionZSTD
	default:
		return sarama.CompressionNone
	}
}

// PublishUserCreated 发送 UserCreated 事件
func (kp *KafkaProducer) PublishUserCreated(event *UserCreated) error {
	// 序列化事件为 JSON
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 构建 Sarama 消息对象
	msg := &sarama.ProducerMessage{
		Topic: event.Topic(),
		Value: sarama.StringEncoder(value),
	}

	// 同步发送消息（阻塞直至成功或失败）
	_, _, err = kp.producer.SendMessage(msg)
	return err
}

func (kp *KafkaProducer) PublishWithRetry(event *UserCreated, maxRetries int) error {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		err := kp.PublishUserCreated(event)
		if err == nil {
			return nil
		}
		lastErr = err
		if i < maxRetries {
			backoff := time.Second * time.Duration(1<<i) // 指数退避：1s, 2s, 4s...
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("failed to publish event after %d retries: %w", maxRetries, lastErr)
}

// Close 关闭生产者并释放资源
func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
