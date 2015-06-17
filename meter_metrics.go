package gorelic

import "github.com/courtf/newrelic_platform_go"

func GetMeterMetrica(ds DataSource, dataSourceKey, basePath, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		NewMeterMetrica(ds, dataSourceKey, basePath, "Count", units, MeterCount),
		NewMeterMetrica(ds, dataSourceKey, basePath, "Rate1", units, MeterRate1),
		NewMeterMetrica(ds, dataSourceKey, basePath, "Rate5", units, MeterRate5),
		NewMeterMetrica(ds, dataSourceKey, basePath, "Rate15", units, MeterRate15),
		NewMeterMetrica(ds, dataSourceKey, basePath, "RateMean", units, MeterRateMean),
	}
}
