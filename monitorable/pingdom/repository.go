package pingdom

import (
	"github.com/monitoror/monitoror/monitorable/pingdom/models"
)

// Repository represent the ping's repository contract
type (
	Repository interface {
		GetCheck(checkId int) (*models.Check, error)
		GetChecks(tags string) ([]models.Check, error)
	}
)
