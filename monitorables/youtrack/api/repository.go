//go:generate mockery --name Repository

package api

type (
	Repository interface {
		CountIssues(query string) (int, error)
	}
)
