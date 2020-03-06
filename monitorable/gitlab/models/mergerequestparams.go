package models

type MergeRequestParams struct {
	Repository string `json:"repository" query:"repository"`
}

func (p *MergeRequestParams) IsValid() bool {
	return p.Repository != ""
}
