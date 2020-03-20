package models

type (
	ChecksParams struct {
		Tags   string `json:"tags" query:"tags"`
		SortBy string `json:"sortBy" query:"sortBy"`
	}
)

func (p *ChecksParams) IsValid() bool {
	if p.SortBy != "" && p.SortBy != "name" {
		return false
	}
	return true
}
