package repository

import (
	"fmt"
	"net/http"
	"time"

	"github.com/monitoror/monitoror/pkg/gopingdom"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	"github.com/monitoror/monitoror/monitorable/pingdom/models"

	. "github.com/jsdidierlaurent/go-pingdom/pingdom"
)

type (
	pingdomRepository struct {
		config *config.Pingdom

		// Pingdom check client
		pingdomCheckApi gopingdom.PingdomCheckApi
	}
)

func NewPingdomRepository(config *config.Pingdom) pingdom.Repository {
	client, err := NewClientWithConfig(ClientConfig{
		BaseURL: config.Url,
		APIKey:  config.ApiKey,
		HTTPClient: &http.Client{
			Timeout: time.Millisecond * time.Duration(config.Timeout),
		},
	})

	// Only if Pingdom Url is not a valid URL
	if err != nil {
		panic(fmt.Sprintf("unable to initiate connection to Pingdom\n. %v\n", err))
	}

	return &pingdomRepository{
		config:          config,
		pingdomCheckApi: client.Checks,
	}
}

func (r *pingdomRepository) GetCheck(id int) (result *models.Check, err error) {
	check, err := r.pingdomCheckApi.Read(id)
	if err != nil {
		return
	}

	result = &models.Check{
		Id:     check.ID,
		Name:   check.Name,
		Status: check.Status,
	}

	return
}

func (r *pingdomRepository) GetChecks(tags string) (results []models.Check, err error) {
	checks, err := r.pingdomCheckApi.List()
	if err != nil {
		return
	}

	for _, check := range checks {
		results = append(results, models.Check{
			Id:     check.ID,
			Name:   check.Name,
			Status: check.Status,
		})
	}

	return
}
