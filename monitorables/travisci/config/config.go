package config

type (
	TravisCI struct {
		URL         string `validate:"required,url,http"`
		Token       string
		GithubToken string
		Timeout     int `validate:"gte=0"` // In Millisecond
	}
)

var Default = &TravisCI{
	URL:         "https://api.travis-ci.com/",
	Token:       "",
	GithubToken: "",
	Timeout:     2000,
}
