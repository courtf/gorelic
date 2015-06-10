package gorelic

import (
	"path/filepath"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

type HistogramMeanMetrica struct {
	baseMetrica
	datasource metrics.Histogram
}

func (m *HistogramMeanMetrica) GetName() string { return m.name }

func (m *HistogramMeanMetrica) GetUnits() string { return m.units }

func (m *HistogramMeanMetrica) GetValue() (float64, error) { return m.datasource.Mean(), nil }

type HistogramMinMetrica struct {
	baseMetrica
	datasource metrics.Histogram
}

func (m *HistogramMinMetrica) GetName() string { return m.name }

func (m *HistogramMinMetrica) GetUnits() string { return m.units }

func (m *HistogramMinMetrica) GetValue() (float64, error) { return float64(m.datasource.Min()), nil }

type HistogramMaxMetrica struct {
	baseMetrica
	datasource metrics.Histogram
}

func (m *HistogramMaxMetrica) GetName() string { return m.name }

func (m *HistogramMaxMetrica) GetUnits() string { return m.units }

func (m *HistogramMaxMetrica) GetValue() (float64, error) { return float64(m.datasource.Max()), nil }

func GetHistogramMetrica(h metrics.Histogram, name, units string) []newrelic_platform_go.IMetrica {
	return []newrelic_platform_go.IMetrica{
		&HistogramMeanMetrica{
			baseMetrica{
				filepath.Join(name, "mean"),
				units,
			},
			h,
		},

		&HistogramMinMetrica{
			baseMetrica{
				filepath.Join(name, "min"),
				units,
			},
			h,
		},

		&HistogramMaxMetrica{
			baseMetrica{
				filepath.Join(name, "max"),
				units,
			},
			h,
		},
	}
}
