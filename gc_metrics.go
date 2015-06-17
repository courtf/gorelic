package gorelic

import (
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

func addGCMericsToComponent(component newrelic_platform_go.IComponent, ds DataSource, pollInterval int) {
	metrics.RegisterDebugGCStats(ds)
	go metrics.CaptureDebugGCStats(ds, time.Duration(pollInterval)*time.Second)

	basePath := "Runtime/GC/"
	component.AddMetrica(NewGaugeMetrica(ds, "debug.GCStats.GCSince", basePath, "Calls", "calls"))
	component.AddMetrica(NewGaugeMetrica(ds, "debug.GCStats.NumGC", basePath, "TotalCalls", "calls"))
	component.AddMetrica(NewGaugeMetrica(ds, "debug.GCStats.PauseTotal", basePath, "PauseTotalTime", "nanos"))

	basePath += "GCTime/"
	hkey := "debug.GCStats.Pause"
	units := "nanos"
	component.AddMetrica(NewHistogramMetrica(ds, hkey, basePath, "Max", units, HistogramMax))
	component.AddMetrica(NewHistogramMetrica(ds, hkey, basePath, "Mean", units, HistogramMean))
	component.AddMetrica(NewHistogramMetrica(ds, hkey, basePath, "Min", units, HistogramMin))
	component.AddMetrica(NewPercentileHistogramMetrica(ds, hkey, basePath, "Percentile95", units, 0.95))
}
