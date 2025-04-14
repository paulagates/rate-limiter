package limiter_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/paulagates/rate-limiter/config"
	"github.com/paulagates/rate-limiter/internal/limiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Get(ctx context.Context, key string) (int, error) {
	args := m.Called(ctx, key)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) Increment(ctx context.Context, key string, ttl time.Duration) (int, error) {
	args := m.Called(ctx, key, ttl)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) SetBlock(ctx context.Context, key string, duration time.Duration) error {
	args := m.Called(ctx, key, duration)
	return args.Error(0)
}

func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) ResetKey(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(1)
}

func TestAllowRequest_TokenWithinLimit(t *testing.T) {
	storage := new(MockStorage)
	cfg := &config.Config{
		RateLimitToken:     5,
		BlockDurationToken: 300,
	}

	l := limiter.NewLimiter(storage, cfg)

	ctx := context.Background()
	token := "test-token"

	storage.On("IsBlocked", ctx, "block:"+token).Return(false, nil)
	storage.On("Increment", ctx, token, time.Second).Return(3, nil)

	allowed, err := l.AllowRequest(ctx, token, true)

	assert.NoError(t, err)
	assert.True(t, allowed)
	storage.AssertExpectations(t)
}

func TestAllowRequest_TokenExceedsLimit(t *testing.T) {
	storage := new(MockStorage)
	cfg := &config.Config{
		RateLimitToken:     2,
		BlockDurationToken: 120,
	}

	l := limiter.NewLimiter(storage, cfg)

	ctx := context.Background()
	token := "spam-token"

	storage.On("IsBlocked", ctx, "block:"+token).Return(false, nil)
	storage.On("Increment", ctx, token, time.Second).Return(3, nil)
	storage.On("SetBlock", ctx, "block:"+token, 120*time.Second).Return(nil)

	allowed, err := l.AllowRequest(ctx, token, true)

	assert.NoError(t, err)
	assert.False(t, allowed)
	storage.AssertExpectations(t)
}

func TestAllowRequest_IPBlocked(t *testing.T) {
	storage := new(MockStorage)
	cfg := &config.Config{
		RateLimitIP:     10,
		BlockDurationIP: 60,
	}

	l := limiter.NewLimiter(storage, cfg)

	ctx := context.Background()
	ip := "192.168.0.99"

	storage.On("IsBlocked", ctx, "block:"+ip).Return(true, nil)

	allowed, err := l.AllowRequest(ctx, ip, false)

	assert.NoError(t, err)
	assert.False(t, allowed)
	storage.AssertExpectations(t)
}

func TestAllowRequest_StorageError(t *testing.T) {
	storage := new(MockStorage)
	cfg := &config.Config{
		RateLimitIP:     10,
		BlockDurationIP: 60,
	}

	l := limiter.NewLimiter(storage, cfg)

	ctx := context.Background()
	ip := "127.0.0.1"

	storage.On("IsBlocked", ctx, "block:"+ip).Return(false, nil)
	storage.On("Increment", ctx, ip, time.Second).Return(0, errors.New("Redis caiu"))

	allowed, err := l.AllowRequest(ctx, ip, false)

	assert.Error(t, err)
	assert.False(t, allowed)
}
