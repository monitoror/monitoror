package models

type (
	TileValue struct {
		Values []string       `json:"values"`
		Unit   TileValuesUnit `json:"unit"`
	}

	TileValuesUnit string
)

const (
	MillisecondUnit TileValuesUnit = "MILLISECOND" // Duration in ms
	RatioUnit       TileValuesUnit = "RATIO"       // Ratio like 0.8465896
	NumberUnit      TileValuesUnit = "NUMBER"      // Number in float
	RawUnit         TileValuesUnit = "RAW"         // String
)

var AvailableTileValuesUnit = map[TileValuesUnit]bool{
	MillisecondUnit: true,
	RatioUnit:       true,
	NumberUnit:      true,
	RawUnit:         true,
}

func (t *Tile) WithValue(unit TileValuesUnit) *Tile {
	t.Value = &TileValue{
		Values: []string{},
		Unit:   unit,
	}
	return t
}
