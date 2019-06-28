package models

type (
	Job struct {
		ID        string
		Buildable bool
		InQueue   bool
	}
)
