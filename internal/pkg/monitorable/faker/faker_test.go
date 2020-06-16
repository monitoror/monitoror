package faker

import (
	"testing"
	"time"

	. "github.com/monitoror/monitoror/models"

	"github.com/stretchr/testify/assert"
)

var status = []Status{
	{SuccessStatus, time.Second * 30},
	{FailedStatus, time.Second * 30},
	{WarningStatus, time.Second * 10},
	{RunningStatus, time.Second * 60},
}

func TestComputeStatus_Panic(t *testing.T) {
	assert.Panics(t, func() { ComputeStatus(time.Now(), Statuses{}) })
}

func TestComputeStatus(t *testing.T) {
	refTime := time.Now()
	assert.Equal(t, SuccessStatus, ComputeStatus(refTime, status))

	// refTime before time.Now()
	assert.Equal(t, SuccessStatus, ComputeStatus(refTime.Add(-time.Second*10), status))
	assert.Equal(t, FailedStatus, ComputeStatus(refTime.Add(-time.Second*40), status))
	assert.Equal(t, WarningStatus, ComputeStatus(refTime.Add(-time.Second*65), status))
	assert.Equal(t, RunningStatus, ComputeStatus(refTime.Add(-time.Second*80), status))
	assert.Equal(t, SuccessStatus, ComputeStatus(refTime.Add(-time.Second*140), status))

	// refTime after time.Now()
	assert.Equal(t, RunningStatus, ComputeStatus(refTime.Add(time.Second*10), status))
	assert.Equal(t, RunningStatus, ComputeStatus(refTime.Add(time.Second*50), status))
	assert.Equal(t, FailedStatus, ComputeStatus(refTime.Add(time.Second*90), status))
}

func BenchmarkComputeStatus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ComputeStatus(time.Now(), status)
	}
}

func TestComputeDuration(t *testing.T) {
	refTime := time.Now()

	assert.InDelta(t, time.Second*30, ComputeDuration(refTime.Add(-time.Second*30), time.Second*300), float64(time.Millisecond*10))
	assert.InDelta(t, time.Second*90, ComputeDuration(refTime.Add(-time.Second*90), time.Second*300), float64(time.Millisecond*10))
	assert.InDelta(t, time.Second*10, ComputeDuration(refTime.Add(-time.Second*310), time.Second*300), float64(time.Millisecond*10))
	assert.InDelta(t, time.Second*10, ComputeDuration(refTime.Add(time.Second*290), time.Second*300), float64(time.Millisecond*10))
}

func BenchmarkComputeDuration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ComputeDuration(time.Now(), time.Second*300)
	}
}

func TestGetRefTime(t *testing.T) {
	assert.True(t, GetRefTime().Before(time.Now()))
}
