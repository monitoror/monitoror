package cache

import (
	"testing"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cache := NewBuildCache(4)

	cache.Add("key", "0", models.SuccessStatus, time.Second*1)
	cache.Add("key", "1", models.SuccessStatus, time.Second*1)
	cache.Add("key", "2", models.SuccessStatus, time.Second*5)
	cache.Add("key", "3", models.SuccessStatus, time.Second*9)
	cache.Add("key", "4", models.SuccessStatus, time.Second*5)
	cache.Add("key", "5", models.FailedStatus, time.Second*1)

	assert.Equal(t, models.FailedStatus, *cache.GetPreviousStatus("key", "6"))
	assert.Equal(t, time.Second*5, *cache.GetEstimatedDuration("key"))
}

func Test_Empty(t *testing.T) {
	cache := NewBuildCache(2)

	assert.Nil(t, cache.GetPreviousStatus("key", "1"))
	assert.Nil(t, cache.GetEstimatedDuration("key"))
}

func Test_AlreadyInCache(t *testing.T) {
	cache := NewBuildCache(4)

	cache.Add("key", "1", models.SuccessStatus, time.Second)
	cache.Add("key", "1", models.SuccessStatus, time.Second)
	assert.Nil(t, cache.GetPreviousStatus("key", "1"))
	assert.Equal(t, time.Second, *cache.GetEstimatedDuration("key"))

	cache.Add("key", "2", models.SuccessStatus, time.Second)
	assert.Equal(t, models.SuccessStatus, *cache.GetPreviousStatus("key", "2"))
}
