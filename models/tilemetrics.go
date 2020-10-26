package models

type (
	TileMetrics struct {
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

func (t *Tile) WithMetrics(unit TileValuesUnit) *Tile {
	t.Metrics = &TileMetrics{
		Values: []string{},
		Unit:   unit,
	}
	return t
}
