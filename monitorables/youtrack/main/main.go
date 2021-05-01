package main

import (
	"fmt"

	"github.com/monitoror/monitoror/monitorables/youtrack/api/repository"
	"github.com/monitoror/monitoror/monitorables/youtrack/config"
)

func main() {
	conf := &config.Youtrack{
		URL:       "http://youtrack.sarbacane.local/",
		Token:     "perm:anNkaWRpZXJsYXVyZW50.NTUtMw==.wWFyzqZpgP4ZcxiBgN6Pg2nEkVGUoQ",
		Timeout:   2000,
		SSLVerify: false,
	}

	repo := repository.NewYoutrackRepository(conf)
	issues, _ := repo.GetIssues("Assignee: rjestin")

	fmt.Println(len(*issues))
}
