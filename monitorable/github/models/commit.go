package models

type (
	Commit struct {
		SHA    string
		Author *Author
	}

	Author struct {
		Name      string
		AvatarURL string
	}
)
