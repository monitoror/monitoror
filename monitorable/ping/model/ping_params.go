//+build !faker

package model

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`
	}
)

func (p *PingParams) IsValid() bool {
	return p.Hostname != ""
}
