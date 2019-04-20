package main

import (
	"context"
	"fmt"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/travisci/repository"
)

func main() {
	conf, _ := config.InitConfig()
	ctx := context.Background()

	repo := repository.NewApiTravisCIRepository(conf)
	build, _ := repo.Build(ctx, "monitoror", "monitoror", "master")

	fmt.Println(build)
}
