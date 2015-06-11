package gorelic

import (
	"fmt"
	"time"

	"github.com/courtf/go-metrics"
)

type HistogramFunc uint8
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

type MetricaDataSource interface {
	metrics.Registry
	GetCounterValue(key string) (float64, error)
	GetGaugeValue(key string) (float64, error)
	GetHistogramValue(key string, hf HistogramFunc, percentile float64) (float64, error)
	GetTimerValue(key string, tf TimerFunc, percentile float64) (float64, error)
	IncCounterForKey(key string, i int64)
	UpdateHistogramForKey(key string, i int64)
	UpdateTimerForKey(key string, d time.Duration)
	UpdateTimerSinceForKey(key string, t time.Time)
	TimerFuncForKey(key string, f func())
}

type metricaDataSource struct {
	metrics.Registry
}

func NewMetricaDataSource(r metrics.Registry) MetricaDataSource {
	return metricaDataSource{r}
}

func (ds metricaDataSource) GetCounterValue(key string) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if counter, ok := valueContainer.(metrics.Counter); ok {
		return float64(counter.Count()), nil
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds metricaDataSource) GetGaugeValue(key string) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if gauge, ok := valueContainer.(metrics.Gauge); ok {
		return float64(gauge.Value()), nil
	} else {
		return 0, fmt.Errorf("metrica container has unexpected type: %T\n", valueContainer)
	}
}

func (ds metricaDataSource) GetHistogramValue(key string, hf HistogramFunc, percentile float64) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if histogram, ok := valueContainer.(metrics.Histogram); ok {
		switch hf {
		default:
			return 0, fmt.Errorf("unsupported stat function for histogram: %s\n", hf)
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

func (ds metricaDataSource) GetTimerValue(key string, tf TimerFunc, percentile float64) (float64, error) {
	if valueContainer := ds.Get(key); valueContainer == nil {
		return 0, fmt.Errorf("metrica with name %s is not registered\n", key)
	} else if timer, ok := valueContainer.(metrics.Timer); ok {
		switch tf {
		default:
			return 0, fmt.Errorf("unsupported stat function for timer: %s\n", tf)
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

func (ds metricaDataSource) counterForKey(key string) (counter metrics.Counter) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	counter, _ = container.(metrics.Counter)
	return
}

func (ds metricaDataSource) histogramForKey(key string) (histogram metrics.Histogram) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	histogram, _ = container.(metrics.Histogram)
	return
}

func (ds metricaDataSource) timerForKey(key string) (timer metrics.Timer) {
	var container interface{}
	if container = ds.Get(key); container == nil {
		return
	}

	timer, _ = container.(metrics.Timer)
	return
}

func (ds metricaDataSource) IncCounterForKey(key string, i int64) {
	if counter := ds.counterForKey(key); counter != nil {
		counter.Inc(i)
	}
}

func (ds metricaDataSource) UpdateHistogramForKey(key string, i int64) {
	if histogram := ds.histogramForKey(key); histogram != nil {
		histogram.Update(i)
	}
}

func (ds metricaDataSource) UpdateTimerForKey(key string, d time.Duration) {
	if timer := ds.timerForKey(key); timer != nil {
		timer.Update(d)
	}
}

func (ds metricaDataSource) UpdateTimerSinceForKey(key string, t time.Time) {
	if timer := ds.timerForKey(key); timer != nil {
		timer.UpdateSince(t)
	}
}

func (ds metricaDataSource) TimerFuncForKey(key string, f func()) {
	if timer := ds.timerForKey(key); timer != nil {
		timer.Time(f)
	}
}
