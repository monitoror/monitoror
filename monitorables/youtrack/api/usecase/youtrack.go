package usecase

import "github.com/monitoror/monitoror/monitorables/youtrack/api"

type (
	youtrackUsecase struct {
		repository api.Repository
	}
)

func NewYoutrackUsecase(repository api.Repository) api.Usecase {
	return &youtrackUsecase{repository}
}
