package main

import (
	"fmt"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/repository"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/usecase"
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
	uc := usecase.NewYoutrackUsecase(repo)

	params := &models.IssuesCountParams{
		Query: "Assignee: rjestin",
		CountThreshold: map[coreModels.TileStatus]int{
			coreModels.FailedStatus: 200,
		},
		PriorityFieldThreshold: map[coreModels.TileStatus]int{
			coreModels.FailedStatus: 1,
		},
	}

	tile, err := uc.CountIssues(params)

	if err != nil {
		panic(err)
	}

	fmt.Println(tile.Metrics.Values)
	fmt.Println(tile.Status)
}
