// internal/messaging/event.go
package messaging

import (
	"time"
)

// UserCreated 表示用户创建成功的事件
type UserCreated struct {
	UserID    string    `json:"user_id"`    // 唯一用户标识
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱地址
	CreatedAt time.Time `json:"created_at"` // 创建时间，RFC3339 格式
}

// Topic 返回该事件对应的主题名称
// 与 config.yaml 中 kafka.topics.user_created 对应
func (e *UserCreated) Topic() string {
	return "user.created.event"
}
