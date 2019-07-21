package cache

import (
	"time"

	"github.com/monitoror/monitoror/models/tiles"
)

type BuildCache struct {
	maxSize        int
	previousBuilds map[string][]build
}

type build struct {
	id       string
	status   tiles.TileStatus
	duration time.Duration
}

func NewBuildCache(size int) *BuildCache {
	return &BuildCache{maxSize: size, previousBuilds: make(map[string][]build)}
}

func (c *BuildCache) GetEstimatedDuration(key string) *time.Duration {
	if _, ok := c.previousBuilds[key]; !ok {
		return nil
	}

	var total int64
	for _, c := range c.previousBuilds[key] {
		total += int64(c.duration)
	}
	average := total / int64(len(c.previousBuilds[key]))
	duration := time.Duration(average)
	return &duration
}

// Get Previous Status excluse current status in case of multiple call with the same current build
func (c *BuildCache) GetPreviousStatus(key, id string) *tiles.TileStatus {
	builds, ok := c.previousBuilds[key]
	if !ok {
		return nil
	}

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
	if _, ok := c.previousBuilds[key]; !ok {
		c.previousBuilds[key] = []build{}
	}

	// if id already exist, skip
	for _, value := range c.previousBuilds[key] {
		if value.id == id {
			return
		}
	}

	// Remove old elements
	if len(c.previousBuilds[key]) == c.maxSize {
		c.previousBuilds[key] = c.previousBuilds[key][:len(c.previousBuilds[key])-1]
	}

	c.previousBuilds[key] = append([]build{{id, s, d}}, c.previousBuilds[key]...)
}
