package models

type PullRequestParams struct {
	Owner      string `json:"owner" query:"owner"`
	Repository string `json:"repository" query:"repository"`
}

func (p *PullRequestParams) IsValid() bool {
	return p.Owner != "" && p.Repository != ""
}
