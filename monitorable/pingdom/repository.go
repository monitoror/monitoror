package pingdom

import (
	"github.com/monitoror/monitoror/monitorable/pingdom/models"
)

type (
	Repository interface {
		GetCheck(checkID int) (*models.Check, error)
		GetChecks(tags string) ([]models.Check, error)
	}
)
