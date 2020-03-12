package stripe

import "github.com/monitoror/monitoror/models"

const (
	StripeCountTileType models.TileType = "STRIPE-COUNT"
)

type (
	Usecase interface {
		Count() (*models.Tile, error)
	}
)
