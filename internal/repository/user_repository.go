// internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"my-go-microservice/internal/model"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

// userRepository implements the UserRepository interface using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID retrieves a user by their ID from the database
func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Return nil, nil to indicate not found
	}
	return &user, err
}

// Update modifies an existing user in the database
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete removes a user from the database by ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
