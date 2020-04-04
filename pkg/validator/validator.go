//go:generate mockery -name SimpleValidator

package validator

type SimpleValidator interface {
	IsValid() bool
}
