package port

// Repository represent the port's repository contract
type (
	Repository interface {
		CheckPort(hostname string, port int) error
	}
)
