package gorelic

import (
	"fmt"
	"time"

	"github.com/courtf/go-metrics"
)

type HistogramFunc uint8
type MeterFunc uint8
type TimerFunc uint8

const (
	HistogramCount HistogramFunc = iota
	HistogramMax
	HistogramMean
	HistogramMin
	HistogramPercentile
	HistogramStdDev
	HistogramSum
	HistogramVariance
	NoHistogramFuncs
)

const (
	MeterCount MeterFunc = iota
	MeterRate1
	MeterRate5
	MeterRate15
	MeterRateMean
)

const (
	TimerCount TimerFunc = iota
	TimerMax
	TimerMean
	TimerMin
	TimerPercentile
	TimerRate1
	TimerRate5
	TimerRate15
	TimerRateMean
	TimerStdDev
	TimerSum
	TimerVariance
	NoTimerFuncs
)

type DataSource interface {
	metrics.Registry
	GetCounterValue(key string) (float64, error)
	GetGaugeValue(key string) (float64, error)
	GetHistogramValue(key string, hf HistogramFunc, percentile float64) (float64, error)
	GetMeterValue(key string, mf MeterFunc) (float64, error)
	GetTimerValue(key string, tf TimerFunc, percentile float64) (float64, error)
	UpdateGaugeForKey(key string, i int64)
	IncCounterForKey(key string, i int64)
	UpdateHistogramForKey(key string, i int64)
	MarkMeterForKey(key string, i int64)
	UpdateTimerForKey(key string, d time.Duration)
	UpdateTimerSinceForKey(key string, t time.Time)
	TimerFuncForKey(key string, f func())
}

type dataSource struct {
	metrics.Registry
}

func NewDataSource(r metrics.Registry) DataSource {
	return dataSource{r}
}

func (ds dataSource) GetCounterValue(key string) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if counter, ok := valueContainer.(metrics.Counter); ok {
		return float64(counter.Count()), nil
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds dataSource) GetGaugeValue(key string) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if gauge, ok := valueContainer.(metrics.Gauge); ok {
		return float64(gauge.Value()), nil
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds dataSource) GetHistogramValue(key string, hf HistogramFunc, percentile float64) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if histogram, ok := valueContainer.(metrics.Histogram); ok {
		switch hf {
		default:
			return 0, fmt.Errorf("unsupported stat function for histogram: %v\n", hf)
		case HistogramCount:
			return float64(histogram.Count()), nil
		case HistogramMax:
			return float64(histogram.Max()), nil
		case HistogramMean:
			return float64(histogram.Mean()), nil
		case HistogramMin:
			return float64(histogram.Min()), nil
		case HistogramPercentile:
			return float64(histogram.Percentile(percentile)), nil
		case HistogramStdDev:
			return float64(histogram.StdDev()), nil
		case HistogramSum:
			return float64(histogram.Sum()), nil
		case HistogramVariance:
			return float64(histogram.Variance()), nil
		}
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds dataSource) GetMeterValue(key string, mf MeterFunc) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if meter, ok := valueContainer.(metrics.Meter); ok {
		switch mf {
		default:
			return 0, fmt.Errorf("unsupported stat function for meter: %v\n", mf)
		case MeterCount:
			return float64(meter.Count()), nil
		case MeterRate1:
			return float64(meter.Rate1()), nil
		case MeterRate5:
			return float64(meter.Rate5()), nil
		case MeterRate15:
			return float64(meter.Rate15()), nil
		case MeterRateMean:
			return float64(meter.RateMean()), nil
		}
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds dataSource) GetTimerValue(key string, tf TimerFunc, percentile float64) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if timer, ok := valueContainer.(metrics.Timer); ok {
		switch tf {
		default:
			return 0, fmt.Errorf("unsupported stat function for timer: %v\n", tf)
		case TimerCount:
			return float64(timer.Count()), nil
		case TimerMax:
			return float64(timer.Max()) / float64(time.Millisecond), nil
		case TimerMean:
			return float64(timer.Mean()) / float64(time.Millisecond), nil
		case TimerMin:
			return float64(timer.Min()) / float64(time.Millisecond), nil
		case TimerPercentile:
			return float64(timer.Percentile(percentile)) / float64(time.Millisecond), nil
		case TimerRate1:
			return float64(timer.Rate1()), nil
		case TimerRate5:
			return float64(timer.Rate5()), nil
		case TimerRate15:
			return float64(timer.Rate15()), nil
		case TimerRateMean:
			return float64(timer.RateMean()), nil
		case TimerStdDev:
			return float64(timer.StdDev()), nil
		case TimerSum:
			return float64(timer.Sum()), nil
		case TimerVariance:
			return float64(timer.Variance()), nil
		}
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds dataSource) gaugeForKey(key string) (gauge metrics.Gauge) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	gauge, _ = container.(metrics.Gauge)
	return
}

func (ds dataSource) counterForKey(key string) (counter metrics.Counter) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	counter, _ = container.(metrics.Counter)
	return
}

func (ds dataSource) histogramForKey(key string) (histogram metrics.Histogram) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	histogram, _ = container.(metrics.Histogram)
	return
}

func (ds dataSource) meterForKey(key string) (meter metrics.Meter) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	meter, _ = container.(metrics.Meter)
	return
}

func (ds dataSource) timerForKey(key string) (timer metrics.Timer) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	timer, _ = container.(metrics.Timer)
	return
}

func (ds dataSource) UpdateGaugeForKey(key string, i int64) {
	if gauge := ds.gaugeForKey(key); gauge != nil {
		gauge.Update(i)
	}
}

func (ds dataSource) IncCounterForKey(key string, i int64) {
	if counter := ds.counterForKey(key); counter != nil {
		counter.Inc(i)
	}
}

func (ds dataSource) UpdateHistogramForKey(key string, i int64) {
	if histogram := ds.histogramForKey(key); histogram != nil {
		histogram.Update(i)
	}
}

func (ds dataSource) MarkMeterForKey(key string, i int64) {
	if meter := ds.meterForKey(key); meter != nil {
		meter.Mark(i)
	}
}

func (ds dataSource) UpdateTimerForKey(key string, d time.Duration) {
	if timer := ds.timerForKey(key); timer != nil {
		timer.Update(d)
	}
}

func (ds dataSource) UpdateTimerSinceForKey(key string, t time.Time) {
	if timer := ds.timerForKey(key); timer != nil {
		timer.UpdateSince(t)
	}
}

func (ds dataSource) TimerFuncForKey(key string, f func()) {
	if timer := ds.timerForKey(key); timer != nil {
		timer.Time(f)
	}
}
