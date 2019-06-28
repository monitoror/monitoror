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

func (c *BuildCache) GetPreviousStatus(key string) *tiles.TileStatus {
	if _, ok := c.previousBuilds[key]; !ok {
		return nil
	}

	return &c.previousBuilds[key][0].status
}

func (c *BuildCache) Add(key string, s tiles.TileStatus, d time.Duration) {
	// If cache is not found, create it
	if _, ok := c.previousBuilds[key]; !ok {
		c.previousBuilds[key] = []build{}
	}

	// Remove old elements
	if len(c.previousBuilds[key]) == c.maxSize {
		c.previousBuilds[key] = c.previousBuilds[key][:len(c.previousBuilds[key])-1]
	}

	c.previousBuilds[key] = append([]build{{s, d}}, c.previousBuilds[key]...)
}
