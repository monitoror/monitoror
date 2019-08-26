//+build !faker

package models

import "regexp"

type (
	MultiBranchParams struct {
		Job    string `json:"job" query:"job"`
		Filter string `json:"filter" query:"filter"`
	}
)

func (p *MultiBranchParams) IsValid() bool {
	if p.Job == "" {
		return false
	}

	if p.Filter != "" {
		_, err := regexp.Compile(p.Filter)
		if err != nil {
			return false
		}
	}

	return true
}
