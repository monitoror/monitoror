//+build !faker

package models

type (
	CheckParams struct {
		ID *int `json:"id" query:"id"`
	}
)

func (p *CheckParams) IsValid() bool {
	return p.ID != nil
}
