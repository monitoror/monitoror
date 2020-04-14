package config

type (
	TravisCI struct {
		URL         string
		Timeout     int // In Millisecond
		Token       string
		GithubToken string
	}
)

var Default = &TravisCI{
	URL:         "https://api.travis-ci.com/",
	Timeout:     2000,
	Token:       "",
	GithubToken: "",
}
