package models

import "regexp"

type (
	MultiBranchParams struct {
		Job string `json:"job" query:"job"`

		// Using Match / Unmatch filter instead of one filter because Golang's standard regex library doesn't have negative look ahead.
		Match   string `json:"match" query:"match"`
		Unmatch string `json:"unmatch" query:"unmatch"`
	}
)

func (p *MultiBranchParams) IsValid() bool {
	if p.Job == "" {
		return false
	}

	if p.Match != "" {
		_, err := regexp.Compile(p.Match)
		if err != nil {
			return false
		}
	}

	if p.Unmatch != "" {
		_, err := regexp.Compile(p.Unmatch)
		if err != nil {
			return false
		}
	}

	return true
}
