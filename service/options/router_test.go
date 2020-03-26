package options

import (
	"testing"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestApplyOptions(t *testing.T) {
	settings := ApplyOptions()
	assert.Len(t, settings.Middlewares, 0)
	assert.Nil(t, settings.CustomCacheExpiration)
	assert.False(t, settings.NoCache)
}

func TestWithMiddlewares(t *testing.T) {
	option := WithMiddlewares(middleware.AddTrailingSlash())
	settings := ApplyOptions(option)
	assert.Len(t, settings.Middlewares, 1)
}

func TestWithCustomCacheExpiration(t *testing.T) {
	option := WithCustomCacheExpiration(time.Second * 5)
	settings := ApplyOptions(option)
	assert.Equal(t, time.Second*5, *settings.CustomCacheExpiration)
}

func TestWithNoCache(t *testing.T) {
	option := WithNoCache()
	settings := ApplyOptions(option)
	assert.True(t, settings.NoCache)
}
