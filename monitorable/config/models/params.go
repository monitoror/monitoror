package models

type (
	ConfigParams struct {
		URL  string `json:"url" query:"url"`
		Path string `json:"path" query:"path"`
	}
)

func (p *ConfigParams) IsValid() bool {
	count := 0
	if p.URL != "" {
		count++
	}
	if p.Path != "" {
		count++
	}
	return count == 1
}
