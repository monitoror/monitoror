package models

type IntThreshold map[TileStatus]int

func (t IntThreshold) GetTileStatus(value int, defaultStatus TileStatus) TileStatus {
	var maxThreshold int
	var thresholdTileType TileStatus

	for tileType, threshold := range t {
		if value > threshold && threshold > maxThreshold {
			maxThreshold = threshold
			thresholdTileType = tileType
		}
	}

	if maxThreshold != 0 {
		return thresholdTileType
	}

	return defaultStatus
}
