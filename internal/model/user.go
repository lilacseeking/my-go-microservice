// internal/model/user.go
package model

import "time"

// User represents the user entity in the system
type User struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"` // Primary key, UUID
	Name      string    `json:"name" gorm:"column:name;not null"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the database table name for GORM
func (User) TableName() string {
	return "users"
}
