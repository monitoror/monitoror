//+build !faker

package model

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`
	}
)

func (p *PingParams) Validate() bool {
	return p.Hostname != ""
}
