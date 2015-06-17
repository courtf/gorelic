package gorelic

import (
	"path/filepath"

	"github.com/courtf/go-metrics"
)

type baseMetrica struct {
	dataSource    DataSource
	dataSourceKey string
	basePath      string
	name          string
	units         string
}

func (metrica baseMetrica) GetName() string {
	return filepath.Join(metrica.basePath, metrica.name)
}

func (metrica baseMetrica) GetUnits() string {
	return metrica.units
}

func (metrica baseMetrica) ClearSentData() {
	// implemented by children (or not)
}

type CounterMetrica struct {
	baseMetrica
}

func NewCounterMetrica(ds DataSource, dataSourceKey, basePath, name, units string) CounterMetrica {
	return CounterMetrica{
		baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
	}
}

func (metrica CounterMetrica) GetValue() (float64, error) {
	return metrica.dataSource.GetCounterValue(metrica.dataSourceKey)
}

func (metrica CounterMetrica) ClearSentData() {
	var container interface{}
	if container = metrica.dataSource.Get(metrica.dataSourceKey); container == nil {
		return
	}

	var counter metrics.Counter
	var ok bool
	if counter, ok = container.(metrics.Counter); !ok {
		return
	}

	counter.Clear()
}

type GaugeMetrica struct {
	baseMetrica
}

func NewGaugeMetrica(ds DataSource, dataSourceKey, basePath, name, units string) GaugeMetrica {
	return GaugeMetrica{
		baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
	}
}

func (metrica GaugeMetrica) GetValue() (float64, error) {
	return metrica.dataSource.GetGaugeValue(metrica.dataSourceKey)
}

type GaugeDeltaMetrica struct {
	baseMetrica
	previousValue float64
}

func NewGaugeDeltaMetrica(ds DataSource, dataSourceKey, basePath, name, units string) *GaugeDeltaMetrica {
	return &GaugeDeltaMetrica{
		baseMetrica: baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
	}
}

func (metrica *GaugeDeltaMetrica) GetValue() (float64, error) {
	var value float64
	var currentValue float64
	var err error
	if currentValue, err = metrica.dataSource.GetGaugeValue(metrica.dataSourceKey); err == nil {
		value = currentValue - metrica.previousValue
		metrica.previousValue = currentValue
	}
	return value, err
}

type HistogramMetrica struct {
	baseMetrica
	histFunc   HistogramFunc
	percentile float64
}

func NewHistogramMetrica(ds DataSource, dataSourceKey, basePath, name, units string,
	hf HistogramFunc) HistogramMetrica {
	return HistogramMetrica{
		baseMetrica: baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
		histFunc: hf,
	}
}

func NewPercentileHistogramMetrica(ds DataSource, dataSourceKey, basePath, name, units string,
	percentile float64) HistogramMetrica {
	return HistogramMetrica{
		baseMetrica: baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
		histFunc:   HistogramPercentile,
		percentile: percentile,
	}
}

func (metrica HistogramMetrica) GetValue() (float64, error) {
	return metrica.dataSource.GetHistogramValue(metrica.dataSourceKey, metrica.histFunc, metrica.percentile)
}

type TimerMetrica struct {
	baseMetrica
	timerFunc  TimerFunc
	percentile float64
}

func NewTimerMetrica(ds DataSource, dataSourceKey, basePath, name, units string, tf TimerFunc) TimerMetrica {
	return TimerMetrica{
		baseMetrica: baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
		timerFunc: tf,
	}
}

func NewPercentileTimerMetrica(ds DataSource, dataSourceKey, basePath, name, units string,
	percentile float64) TimerMetrica {
	return TimerMetrica{
		baseMetrica: baseMetrica{
			ds, dataSourceKey, basePath, name, units,
		},
		timerFunc:  TimerPercentile,
		percentile: percentile,
	}
}

func (metrica TimerMetrica) GetValue() (float64, error) {
	return metrica.dataSource.GetTimerValue(metrica.dataSourceKey, metrica.timerFunc, metrica.percentile)
}
