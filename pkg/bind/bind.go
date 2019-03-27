package bind

type (
	Binder interface {
		// Bind information into parameters pointer
		Bind(i interface{}) error
	}
)
