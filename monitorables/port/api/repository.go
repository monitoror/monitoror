package api

type (
	Repository interface {
		OpenSocket(hostname string, port int) error
	}
)
