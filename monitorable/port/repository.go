package port

// Repository represent the port's repository contract
type (
	Repository interface {
		OpenSocket(hostname string, port int) error
	}
)
