//+build faker

package monitorable

import coreModels "github.com/monitoror/monitoror/models"

type DefaultMonitorableFaker struct {
}

func (m *DefaultMonitorableFaker) GetVariants() []coreModels.VariantName {
	return []coreModels.VariantName{coreModels.DefaultVariant}
}

func (m *DefaultMonitorableFaker) Validate(_ coreModels.VariantName) (bool, error) {
	return true, nil
}
