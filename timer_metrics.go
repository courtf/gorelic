package gorelic

import "github.com/courtf/newrelic_platform_go"

func addTimerMeterMetrics(component newrelic_platform_go.IComponent, ds DataSource, dataSourceKey, basePath, units string) {
	for _, m := range GetTimerMeterMetrica(ds, dataSourceKey, basePath, units) {
		component.AddMetrica(m)
	}
}

func GetTimerMeterMetrica(ds DataSource, dataSourceKey, basePath, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		NewTimerMetrica(ds, dataSourceKey, basePath, "Rate1", units, TimerRate1),
		NewTimerMetrica(ds, dataSourceKey, basePath, "Rate5", units, TimerRate5),
		NewTimerMetrica(ds, dataSourceKey, basePath, "Rate15", units, TimerRate15),
		NewTimerMetrica(ds, dataSourceKey, basePath, "RateMean", units, TimerRateMean),
	}
}

func addTimerHistogramMetrics(component newrelic_platform_go.IComponent, ds DataSource, dataSourceKey, basePath string) {
	for _, m := range GetTimerHistogramMetrica(ds, dataSourceKey, basePath) {
		component.AddMetrica(m)
	}
}

func GetTimerHistogramMetrica(ds DataSource, dataSourceKey, basePath string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		NewTimerMetrica(ds, dataSourceKey, basePath, "Max", "ms", TimerMax),
		NewTimerMetrica(ds, dataSourceKey, basePath, "Mean", "ms", TimerMean),
		NewTimerMetrica(ds, dataSourceKey, basePath, "Min", "ms", TimerMin),
		NewPercentileTimerMetrica(ds, dataSourceKey, basePath, "Percentile95", "ms", 0.95),
	}
}

func GetTimerMetrica(ds DataSource, dataSourceKey, basePath, units string) []newrelic_platform_go.IMetrica {
	mm := GetTimerMeterMetrica(ds, dataSourceKey, basePath, units)
	mmLen := len(mm)
	thm := GetTimerHistogramMetrica(ds, dataSourceKey, basePath)
	thmLen := len(thm)
	total := mmLen + thmLen

	ret := make([]newrelic_platform_go.IMetrica, 0, total)
	for i := 0; i < total; i++ {
		if i < mmLen {
			ret[i] = mm[i]
			continue
		}

		ret[i] = thm[i]
	}

	return ret
}
