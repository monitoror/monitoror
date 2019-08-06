package cache

import (
	"time"

	cmap "github.com/orcaman/concurrent-map"

	"github.com/monitoror/monitoror/models/tiles"
)

type BuildCache struct {
	maxSize        int
	previousBuilds cmap.ConcurrentMap
}

type build struct {
	id       string
	status   tiles.TileStatus
	duration time.Duration
}

func NewBuildCache(size int) *BuildCache {
	return &BuildCache{maxSize: size, previousBuilds: cmap.New()}
}

func (c *BuildCache) GetEstimatedDuration(key string) *time.Duration {
	value, ok := c.previousBuilds.Get(key)
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

// Get Previous Status excluse current status in case of multiple call with the same current build
func (c *BuildCache) GetPreviousStatus(key, id string) *tiles.TileStatus {
	value, ok := c.previousBuilds.Get(key)
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

func (c *BuildCache) Add(key, id string, s tiles.TileStatus, d time.Duration) {
	// If cache is not found, create it
	var builds []build
	if tmp, ok := c.previousBuilds.Get(key); !ok {
		c.previousBuilds.Set(key, builds)
	} else {
		builds = tmp.([]build)
	}

	// if id already exist, skip
	for _, value := range builds {
		if value.id == id {
			return
		}
	}

	// Remove old elements
	if len(builds) == c.maxSize {
		builds = builds[:len(builds)-1]
	}

	c.previousBuilds.Set(key, append([]build{{id, s, d}}, builds...))
}
