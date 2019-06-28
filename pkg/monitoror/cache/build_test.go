package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/models/tiles"
)

func Test(t *testing.T) {
	cache := NewBuildCache(2)

	cache.Add("key", tiles.SuccessStatus, time.Second*9)
	cache.Add("key", tiles.SuccessStatus, time.Second*9)
	cache.Add("key", tiles.SuccessStatus, time.Second*9)
	cache.Add("key", tiles.SuccessStatus, time.Second*9)
	cache.Add("key", tiles.SuccessStatus, time.Second*9)
	cache.Add("key", tiles.FailedStatus, time.Second*1)

	assert.Equal(t, tiles.FailedStatus, *cache.GetPreviousStatus("key"))
	assert.Equal(t, time.Second*5, *cache.GetEstimatedDuration("key"))
}

func Test_Empty(t *testing.T) {
	cache := NewBuildCache(2)

	assert.Nil(t, cache.GetPreviousStatus("key"))
	assert.Nil(t, cache.GetEstimatedDuration("key"))
}
