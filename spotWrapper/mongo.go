package spotwrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// cache interface defining methods for caching
type cache[T comparable, V any] interface {
	Get(ctx context.Context, key string) V
	Set(ctx context.Context, key string, data T, expire int) // Handles inserts & updates
	Delete(ctx context.Context, key string)
	Exist(ctx context.Context, key string) bool
}

// Factory function to create a new cache instance
func NewCache[T comparable, V ~string]() cache[T, V] {
	return newSpotCache[T, V]()
}

// spotish struct implementing the cache interface
type spotish[T comparable, V any] struct {
	client *redis.Client
}

// Function to initialize a new Redis client
func newSpotCache[T comparable, V any]() *spotish[T, V] {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	return &spotish[T, V]{client: client}
}

// Get method retrieves a value from Redis cache
func (s *spotish[T, V]) Get(ctx context.Context, key string) V {
	str, err := s.client.Get(key).Result()
	if err != nil {
		var zero V
		return zero
	}

	// Convert string back to type V
	value, ok := any(str).(V)
	if !ok {
		var zero V
		return zero
	}

	return value
}

// Set method stores a value in Redis cache
func (s *spotish[T, V]) Set(ctx context.Context, key string, data T, expire int) {
	err := s.client.Set(key, fmt.Sprintf("%v", data), time.Duration(expire)*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting cache:", err)
	}
}

// Delete method removes a key from Redis cache
func (s *spotish[T, V]) Delete(ctx context.Context, key string) {
	err := s.client.Del(key).Err()
	if err != nil {
		fmt.Println("Error deleting cache key:", err)
	}
}

// Exist method checks if a key exists in Redis cache
func (s *spotish[T, V]) Exist(ctx context.Context, key string) bool {
	exists, err := s.client.Exists(key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

func SaveUser() {
	fmt.Println("Prentending to save user info ->  implement later")
}
