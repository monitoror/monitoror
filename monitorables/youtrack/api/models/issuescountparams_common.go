package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/monitoror/monitoror/internal/pkg/validator"
	coreModels "github.com/monitoror/monitoror/models"
)

const DefaultPriorityFieldLabel = "Priority"
const DefaultPriorityFieldValue = "Show-stopper"

func buildPriorityQuery(query string, label, value *string) string {
	fieldLabel := DefaultPriorityFieldLabel
	if label != nil {
		fieldLabel = *label
	}
	fieldValue := DefaultPriorityFieldValue
	if value != nil {
		fieldValue = *value
	}

	// YouTrack syntax requires braces around field names containing spaces
	if strings.ContainsAny(fieldLabel, " \t") {
		fieldLabel = fmt.Sprintf("{%s}", fieldLabel)
	}

	return fmt.Sprintf("(%s) and %s: {%s}", query, fieldLabel, fieldValue)
}

func validateThresholds(countThreshold, priorityFieldThreshold coreModels.IntThreshold) []validator.Error {
	var errs []validator.Error
	errs = append(errs, validateThreshold("CountThreshold", countThreshold)...)
	errs = append(errs, validateThreshold("PriorityFieldThreshold", priorityFieldThreshold)...)
	return errs
}

func validateThreshold(name string, threshold coreModels.IntThreshold) []validator.Error {
	var errs []validator.Error
	for tileStatus := range threshold {
		if !slices.Contains(coreModels.ThresholdAuthorizedTileStatus, tileStatus) {
			errs = append(errs, validator.NewDefaultError(fmt.Sprintf("%s.%s", name, string(tileStatus)), fmt.Sprintf(" %v", coreModels.ThresholdAuthorizedTileStatus)))
		}
	}
	return errs
}
