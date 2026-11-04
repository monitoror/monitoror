package models

type (
	Monitor struct {
		ID   string
		Name string

		IsDisabled    bool
		IsOperational bool
		IsOffline     bool

		Labels []string
	}
)

func (m Monitor) MatchOneLabel(labels []string) bool {
	for _, l := range m.Labels {
		for _, label := range labels {
			if l == label {
				return true
			}
		}
	}
	return false
}
