package spotwrapper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCache checks the basic functionality of the cache implementation
func TestCache(t *testing.T) {
	ctx := context.Background()
	cache := NewCache[string, string]() // Initialize cache with string keys & values

	// Test Set()
	cache.Set(ctx, "testKey", "testValue", 5) // Expires in 5 sec

	// Test Get()
	val := cache.Get(ctx, "testKey")
	assert.Equal(t, "testValue", val, "Expected value should match stored value")

	// Test Exist()
	exists := cache.Exist(ctx, "testKey")
	assert.True(t, exists, "Key should exist after being set")

	// Test Delete()
	cache.Delete(ctx, "testKey")
	existsAfterDelete := cache.Exist(ctx, "testKey")
	assert.False(t, existsAfterDelete, "Key should not exist after deletion")

	// Test Get() on missing key
	missingVal := cache.Get(ctx, "nonExistingKey")
	assert.Equal(t, "", missingVal, "Get() should return empty string for missing keys")

	// Wait for key expiration and test existence
	time.Sleep(6 * time.Second)
	expiredVal := cache.Get(ctx, "testKey")
	assert.Equal(t, "", expiredVal, "Key should have expired after timeout")
}
