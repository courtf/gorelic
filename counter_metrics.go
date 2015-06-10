package gorelic

import (
	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

type baseMetrica struct {
	name, units string
}

type CounterMetrica struct {
	baseMetrica
	datasource metrics.Counter
}

func (m *CounterMetrica) GetName() string { return m.name }

func (m *CounterMetrica) GetUnits() string { return m.units }

func (m *CounterMetrica) GetValue() (float64, error) { return float64(m.datasource.Count()), nil }

func GetCounterMetrica(c metrics.Counter, name, units string) newrelic_platform_go.IMetrica {
	return &CounterMetrica{
		baseMetrica{
			name:  name,
			units: units,
		},
		c,
	}
}
