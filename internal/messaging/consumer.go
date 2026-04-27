// internal/messaging/consumer.go
package messaging

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"my-go-microservice/internal/config"
)

// KafkaConsumer 封装 Sarama 消费者组客户端
type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
	handler       MessageHandler
}

// MessageHandler 定义消息处理回调接口
type MessageHandler interface {
	HandleMessage(event *UserCreated) error
}

// NewKafkaConsumer 根据配置创建并启动 Kafka 消费者组
func NewKafkaConsumer(cfg *config.Kafka, handler MessageHandler) (*KafkaConsumer, error) {
	// 创建 Sarama 配置实例
	config := sarama.NewConfig()

	// 设置消费者组会话超时和心跳间隔
	config.Consumer.Group.Session.Timeout = time.Duration(cfg.Consumer.SessionTimeout)
	config.Consumer.Group.Heartbeat.Interval = time.Duration(cfg.Consumer.HeartbeatInterval)

	// 设置初始偏移量策略：oldest 或 newest
	switch cfg.Consumer.InitialOffset {
	case "oldest":
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		fallthrough
	default:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	// 启用自动提交偏移量（生产环境建议关闭，手动控制）
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// 创建消费者组实例
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.Consumer.GroupID, config)
	if err != nil {
		return nil, err
	}

	// 订阅 user_created 主题
	topics := []string{cfg.Topics.UserCreated}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		topics:        topics,
		handler:       handler,
	}, nil
}

// Consume 启动消费者循环，持续拉取消息
func (kc *KafkaConsumer) Consume() error {
	// 持续运行消费循环
	for {
		// 使用 ConsumerGroup.Consume 阻塞等待消息
		err := kc.consumerGroup.Consume(
			nil, // context 可传入以支持取消
			kc.topics,
			kc,
		)
		if err != nil {
			return err
		}

		// 如果返回 nil，表示 Rebalance 发生，继续循环
	}
}

// 实现 sarama.ConsumerGroupHandler 接口
func (kc *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (kc *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 处理单个分区的消息流
func (kc *KafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// 反序列化消息为 UserCreated 事件
		var event UserCreated
		err := json.Unmarshal(message.Value, &event)
		if err != nil {
			// 日志记录后跳过该消息，防止阻塞分区
			continue
		}

		// 调用业务处理器
		err = kc.handler.HandleMessage(&event)
		if err != nil {
			// 处理失败：不提交 Offset，下次重试
			continue
		}

		// 成功处理后提交 Offset
		sess.MarkMessage(message, "")
	}

	return nil
}

// Close 关闭消费者并停止消费循环
func (kc *KafkaConsumer) Close() error {
	return kc.consumerGroup.Close()
}
