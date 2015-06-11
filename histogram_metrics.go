package gorelic

import "github.com/courtf/newrelic_platform_go"

func GetHistogramMetrica(ds DataSource, dataSourceKey, basePath, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		NewHistogramMetrica(ds, dataSourceKey, basePath, "Max", units, HistogramMax),
		NewHistogramMetrica(ds, dataSourceKey, basePath, "Mean", units, HistogramMean),
		NewHistogramMetrica(ds, dataSourceKey, basePath, "Min", units, HistogramMin),
		NewPercentileHistogramMetrica(ds, dataSourceKey, basePath, "Percentile95", units, 0.95),
	}
}
