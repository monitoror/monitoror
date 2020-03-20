//+build faker

package models

import coreModels "github.com/monitoror/monitoror/models"

type (
	CheckParams struct {
		ID *int `json:"id" query:"id"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)

func (p *CheckParams) IsValid() bool {
	return p.ID != nil
}
