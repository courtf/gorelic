package gorelic

import (
	"path/filepath"

	"github.com/courtf/newrelic_platform_go"
)

func GetHistogramMetrica(ds DataSource, dataSourceKey, basePath, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		NewHistogramMetrica(ds, dataSourceKey, filepath.Join(basePath, "Max"), units, HistogramMax),
		NewHistogramMetrica(ds, dataSourceKey, filepath.Join(basePath, "Mean"), units, HistogramMean),
		NewHistogramMetrica(ds, dataSourceKey, filepath.Join(basePath, "Min"), units, HistogramMin),
		NewPercentileHistogramMetrica(ds, dataSourceKey, filepath.Join(basePath, "Percentile95"), units, 0.95),
	}
}
