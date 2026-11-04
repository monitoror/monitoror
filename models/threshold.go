package models

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	// Ordered from the least to the most severe status
	ThresholdAuthorizedTileStatus = []TileStatus{
		DisabledStatus,
		SuccessStatus,
		WarningStatus,
		FailedStatus,
	}
)

type IntThreshold map[TileStatus]int

func (t *IntThreshold) GetTileStatus(value int, defaultStatus TileStatus) TileStatus {
	var maxThreshold int
	var thresholdTileType TileStatus

	for tileType, threshold := range *t {
		if value >= threshold && threshold > maxThreshold {
			maxThreshold = threshold
			thresholdTileType = tileType
		}
	}

	if maxThreshold != 0 {
		return thresholdTileType
	}

	return defaultStatus
}

// MostSevereTileStatus returns the status ranked highest in ThresholdAuthorizedTileStatus
func MostSevereTileStatus(a, b TileStatus) TileStatus {
	if slices.Index(ThresholdAuthorizedTileStatus, b) > slices.Index(ThresholdAuthorizedTileStatus, a) {
		return b
	}
	return a
}

// UnmarshalParam used to parse a string representation of a map inject by queryParam with vue"
func (t *IntThreshold) UnmarshalParam(src string) error {
	// Remove map[ prefix and ] suffix
	src = strings.TrimPrefix(src, "map[")
	src = strings.TrimSuffix(src, "]")
	if src == "" {
		*t = make(map[TileStatus]int)
		return nil
	}

	result := make(map[TileStatus]int)
	entries := strings.Split(src, " ")
	for _, entry := range entries {
		kv := strings.SplitN(strings.TrimSpace(entry), ":", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid format: %s", entry)
		}

		// Check if the value is a valid int
		val, err := strconv.Atoi(kv[1])
		if err != nil {
			return fmt.Errorf("non-numeric value for %s", kv[0])
		}

		// Check if the key is a valid TileStatus
		ts := TileStatus(kv[0])
		if !slices.Contains(ThresholdAuthorizedTileStatus, ts) {
			return fmt.Errorf("invalid TileStatus: %s", kv[0])
		}

		result[ts] = val
	}

	*t = result
	return nil
}
