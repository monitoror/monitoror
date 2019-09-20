package middlewares

import (
	"testing"
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/jsdidierlaurent/echo-middleware/cache/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestNewCacheMiddleware(t *testing.T) {
	store := cache.NewGoCacheStore(time.Second, time.Second)
	middleware := NewCacheMiddleware(store, time.Second, time.Second)

	if assert.NotNil(t, middleware) {
		assert.NotNil(t, middleware.store)
	}
}

func TestUpstreamCacheHandler(t *testing.T) {
	middleware := &CacheMiddleware{store: &upstreamStore{}}
	handle := middleware.UpstreamCacheHandler(func(c echo.Context) error {
		return nil
	})

	assert.NotNil(t, handle)
}

func TestUpstreamCacheHandlerWithExpiration(t *testing.T) {
	middleware := &CacheMiddleware{store: &upstreamStore{}}
	handle := middleware.UpstreamCacheHandlerWithExpiration(time.Hour, func(c echo.Context) error {
		return nil
	})

	assert.NotNil(t, handle)
}

func TestDownstreamStoreMiddleware(t *testing.T) {
	middleware := &CacheMiddleware{store: &upstreamStore{}}
	handle := middleware.DownstreamStoreMiddleware()

	assert.NotNil(t, handle)
}

func TestStore(t *testing.T) {
	mockStore := new(mocks.Store)
	mockStore.On("Get", AnythingOfType("string"), Anything).Return(nil)
	mockStore.On("Set", AnythingOfType("string"), Anything, AnythingOfType("time.Duration")).Return(nil)

	store := &upstreamStore{
		store: mockStore,
	}

	// Test GET
	if assert.NoError(t, store.Get("key", "value")) {
		mockStore.AssertNumberOfCalls(t, "Get", 1)
	}

	// Test SET
	if assert.NoError(t, store.Set("key", "value", time.Hour)) {
		mockStore.AssertNumberOfCalls(t, "Set", 2)
	}

	// Test Add
	assert.Panics(t, func() { _ = store.Add("key", "value", time.Hour) })
	// Test Replace
	assert.Panics(t, func() { _ = store.Replace("key", "value", time.Hour) })
	// Test Delete
	assert.Panics(t, func() { _ = store.Delete("key") })
	// Test Increment
	assert.Panics(t, func() { _, _ = store.Increment("key", uint64(1)) })
	// Test Decrement
	assert.Panics(t, func() { _, _ = store.Decrement("key", uint64(1)) })
	// Test Flush
	assert.Panics(t, func() { _ = store.Flush() })

	mockStore.AssertExpectations(t)
}
