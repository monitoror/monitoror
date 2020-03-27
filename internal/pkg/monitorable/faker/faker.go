package faker

// Only used by faker but i need unit testing it

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/models"
)

type Statuses []Status
type Status struct {
	Status   models.TileStatus
	Duration time.Duration
}

const RefTimeDelta = time.Minute * 5

func (statuses Statuses) GetTotalStatusDuration() time.Duration {
	var totalStatusDuration time.Duration
	for _, status := range statuses {
		totalStatusDuration += status.Duration
	}

	return totalStatusDuration
}

// ComputeStatus determines the tile Status using current time.
func ComputeStatus(refTime time.Time, availableStatus Statuses) models.TileStatus {
	totalStatusDuration := availableStatus.GetTotalStatusDuration()
	if totalStatusDuration == 0 {
		panic("availableStatus can't be empty. Check your code.")
	}

	computedDuration := time.Since(refTime)
	if computedDuration < 0 {
		computedDuration = computedDuration + totalStatusDuration
	}
	computedDuration = computedDuration % totalStatusDuration

	var status models.TileStatus
	for _, fakerStatus := range availableStatus {
		if computedDuration <= fakerStatus.Duration {
			status = fakerStatus.Status
			break
		}
		computedDuration = computedDuration - fakerStatus.Duration
	}

	return status
}

// ComputeStatus determines the tile build duration using current time.
func ComputeDuration(refTime time.Time, duration time.Duration) time.Duration {
	computedDuration := time.Since(refTime)
	if computedDuration < 0 {
		computedDuration = computedDuration + duration
	}
	computedDuration = computedDuration % duration

	return computedDuration
}

// GetRefTime return current time minus random on 5m delta. It's used to randomize staring time of faker
func GetRefTime() time.Time {
	return time.Now().Add(-time.Duration(rand.Int63n(int64(RefTimeDelta))))
}
