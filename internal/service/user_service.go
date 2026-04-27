// internal/service/user_service.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"my-go-microservice/internal/messaging"
	"my-go-microservice/internal/model"
	"my-go-microservice/internal/repository"
)

// UserService defines the interface for user business operations
type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id string) (*model.User, error)
}

// userService implements the UserService interface
type userService struct {
	repo      repository.UserRepository
	redisPool *redis.Pool
	producer  *messaging.KafkaProducer
}

// NewUserService creates a new instance of userService
func NewUserService(repo repository.UserRepository, pool *redis.Pool, producer *messaging.KafkaProducer) UserService {
	return &userService{
		repo:      repo,
		redisPool: pool,
		producer:  producer,
	}
}

// GetUser retrieves a user by ID with Redis cache integration
func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	// Step 1: Generate Redis key
	cacheKey := fmt.Sprintf("user:%s", id)

	// Step 2: Try to get from Redis cache
	conn := s.redisPool.Get()
	defer conn.Close()

	cachedData, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == nil {
		// Cache hit: Unmarshal and return
		var user model.User
		if json.Unmarshal(cachedData, &user) == nil {
			return &user, nil
		}
		// If unmarshal fails, treat as cache miss and proceed to DB
	}

	// Step 3: Cache miss, query database
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // User not found
	}

	// Step 4: Update cache asynchronously
	go func() {
		conn := s.redisPool.Get()
		defer conn.Close()

		data, _ := json.Marshal(user)
		// Set cache with TTL of 5 minutes and handle potential errors silently
		_, _ = conn.Do("SETEX", cacheKey, 300, data)
	}()

	return user, nil
}

// CreateUser creates a new user and publishes a UserCreated event
func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	// Step 1: Persist user to database within a transaction
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	// Step 2: Construct the UserCreated event
	event := &messaging.UserCreated{
		UserID:    user.ID,
		Username:  user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// Step 3: Publish event asynchronously to avoid blocking the main flow
	go func() {
		if err := s.producer.PublishUserCreated(event); err != nil {
			// Log the error but do not fail the main operation
			log.Printf("Failed to publish UserCreated event for user %s: %v", user.ID, err)
		} else {
			log.Printf("Successfully published UserCreated event for user %s", user.ID)
		}
	}()

	return nil
}
