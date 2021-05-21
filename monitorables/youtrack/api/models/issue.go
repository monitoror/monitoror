package models

type (
	Issues []Issue

	Issue struct {
		ID           string        `json:"id"`
		CustomFields []CustomField `json:"customFields"`
	}

	CustomField struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"` // Due to complex dynamic structure, we use interface{} and use reflect.TypeOf in usecase
	}
)
