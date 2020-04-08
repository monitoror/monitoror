//go:generate mockery -name ParamsValidator -output ../mocks

package models

type ParamsValidator interface {
	Validate(currentVersion *ConfigVersion) *ConfigError
}
