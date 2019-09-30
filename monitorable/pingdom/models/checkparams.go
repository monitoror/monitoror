//+build !faker

package models

type (
	CheckParams struct {
		Id *int `json:"id" query:"id"`
	}
)

func (p *CheckParams) IsValid() bool {
	return p.Id != nil
}
