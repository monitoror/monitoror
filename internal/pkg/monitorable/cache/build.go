package cache

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"
	cmap "github.com/orcaman/concurrent-map"
)

type BuildCache struct {
	maxSize        int
	previousBuilds cmap.ConcurrentMap
}

type build struct {
	id       string
	status   models.TileStatus
	duration time.Duration
}

func NewBuildCache(size int) *BuildCache {
	return &BuildCache{maxSize: size, previousBuilds: cmap.New()}
}

func (c *BuildCache) GetEstimatedDuration(key interface{}) *time.Duration {
	k := fmt.Sprint(key)
	value, ok := c.previousBuilds.Get(k)
	if !ok {
		return nil
	}
	builds := value.([]build)

	var total int64
	for _, c := range builds {
		total += int64(c.duration)
	}
	average := total / int64(len(builds))
	duration := time.Duration(average)
	return &duration
}

// Get Previous Status excludes current status in case of multiple call with the same current build
func (c *BuildCache) GetPreviousStatus(key interface{}, id string) *models.TileStatus {
	k := fmt.Sprint(key)
	value, ok := c.previousBuilds.Get(k)
	if !ok {
		return nil
	}
	builds := value.([]build)

	previous := builds[0]
	if previous.id == id {
		if len(builds) == 1 {
			return nil
		}
		previous = builds[1]
	}

	return &previous.status
}

func (c *BuildCache) Add(key interface{}, id string, s models.TileStatus, d time.Duration) {
	k := fmt.Sprint(key)

	// Upsert keeps the read-modify-write atomic and never exposes an empty slice to readers
	c.previousBuilds.Upsert(k, nil, func(exist bool, valueInMap interface{}, _ interface{}) interface{} {
		var builds []build
		if exist {
			builds = valueInMap.([]build)
		}

		// if id already exist, skip
		for _, value := range builds {
			if value.id == id {
				return builds
			}
		}

		// Remove old elements
		if len(builds) == c.maxSize {
			builds = builds[:len(builds)-1]
		}

		return append([]build{{id, s, d}}, builds...)
	})
}
