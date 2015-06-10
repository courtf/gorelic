package gorelic

import (
	"path/filepath"
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

type baseTimerMetrica struct {
	dataSource metrics.Timer
	name       string
	units      string
}

func (metrica *baseTimerMetrica) GetName() string {
	return metrica.name
}

func (metrica *baseTimerMetrica) GetUnits() string {
	return metrica.units
}

func addMeterMetrics(component newrelic_platform_go.IComponent, timer metrics.Timer, name, units string) {
	for _, m := range GetMeterMetrica(timer, name, units) {
		component.AddMetrica(m)
	}
}

func GetMeterMetrica(timer metrics.Timer, name, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		&TimerRate1Metrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "1minute"),
				units:      units,
				dataSource: timer,
			},
		},

		&TimerRate5Metrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "5minute"),
				units:      units,
				dataSource: timer,
			},
		},

		&TimerRate15Metrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "15minute"),
				units:      units,
				dataSource: timer,
			},
		},

		&TimerRateMeanMetrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "rateMean"),
				units:      units,
				dataSource: timer,
			},
		},
	}
}

func addTimedHistogramMetrics(component newrelic_platform_go.IComponent, timer metrics.Timer, name string) {
	for _, m := range GetTimedHistogramMetrica(timer, name) {
		component.AddMetrica(m)
	}
}

func GetTimedHistogramMetrica(timer metrics.Timer, name string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		&TimerMeanMetrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "mean"),
				units:      "ms",
				dataSource: timer,
			},
		},

		&TimerMaxMetrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "max"),
				units:      "ms",
				dataSource: timer,
			},
		},

		&TimerMinMetrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "min"),
				units:      "ms",
				dataSource: timer,
			},
		},

		&TimerPercentile75Metrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "percentile75"),
				units:      "ms",
				dataSource: timer,
			},
		},

		&TimerPercentile90Metrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "percentile90"),
				units:      "ms",
				dataSource: timer,
			},
		},

		&TimerPercentile95Metrica{
			baseTimerMetrica: &baseTimerMetrica{
				name:       filepath.Join(name, "percentile95"),
				units:      "ms",
				dataSource: timer,
			},
		},
	}
}

func GetTimerMetrica(timer metrics.Timer, name, units string) []newrelic_platform_go.IMetrica {
	mm := GetMeterMetrica(timer, name, units)
	mmLen := len(mm)
	thm := GetTimedHistogramMetrica(timer, name)
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

type TimerRate1Metrica struct {
	*baseTimerMetrica
}

func (metrica *TimerRate1Metrica) GetValue() (float64, error) {
	return metrica.dataSource.Rate1(), nil
}

type TimerRate5Metrica struct {
	*baseTimerMetrica
}

func (metrica *TimerRate5Metrica) GetValue() (float64, error) {
	return metrica.dataSource.Rate5(), nil
}

type TimerRate15Metrica struct {
	*baseTimerMetrica
}

func (metrica *TimerRate15Metrica) GetValue() (float64, error) {
	return metrica.dataSource.Rate15(), nil
}

type TimerRateMeanMetrica struct {
	*baseTimerMetrica
}

func (metrica *TimerRateMeanMetrica) GetValue() (float64, error) {
	return metrica.dataSource.RateMean(), nil
}

type TimerMeanMetrica struct {
	*baseTimerMetrica
}

func (metrica *TimerMeanMetrica) GetValue() (float64, error) {
	return metrica.dataSource.Mean() / float64(time.Millisecond), nil
}

type TimerMinMetrica struct {
	*baseTimerMetrica
}

func (metrica *TimerMinMetrica) GetValue() (float64, error) {
	return float64(metrica.dataSource.Min()) / float64(time.Millisecond), nil
}

type TimerMaxMetrica struct {
	*baseTimerMetrica
}

func (metrica *TimerMaxMetrica) GetValue() (float64, error) {
	return float64(metrica.dataSource.Max()) / float64(time.Millisecond), nil
}

type TimerPercentile75Metrica struct {
	*baseTimerMetrica
}

func (metrica *TimerPercentile75Metrica) GetValue() (float64, error) {
	return metrica.dataSource.Percentile(0.75) / float64(time.Millisecond), nil
}

type TimerPercentile90Metrica struct {
	*baseTimerMetrica
}

func (metrica *TimerPercentile90Metrica) GetValue() (float64, error) {
	return metrica.dataSource.Percentile(0.90) / float64(time.Millisecond), nil
}

type TimerPercentile95Metrica struct {
	*baseTimerMetrica
}

func (metrica *TimerPercentile95Metrica) GetValue() (float64, error) {
	return metrica.dataSource.Percentile(0.95) / float64(time.Millisecond), nil
}
