package models

type (
	TileGeneratorFunction func(params interface{}) ([]GeneratedTile, error)

	GeneratedTile struct {
		Label  string
		Params interface{}
	}
)
