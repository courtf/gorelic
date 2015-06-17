package gorelic

import (
	"path/filepath"

	"github.com/courtf/newrelic_platform_go"
)

func GetMeterMetrica(ds DataSource, dataSourceKey, basePath, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		NewMeterMetrica(ds, dataSourceKey, filepath.Join(basePath, "Rate1"), units, MeterRate1),
		NewMeterMetrica(ds, dataSourceKey, filepath.Join(basePath, "Rate5"), units, MeterRate5),
		NewMeterMetrica(ds, dataSourceKey, filepath.Join(basePath, "Rate15"), units, MeterRate15),
		NewMeterMetrica(ds, dataSourceKey, filepath.Join(basePath, "RateMean"), units, MeterRateMean),
	}
}
