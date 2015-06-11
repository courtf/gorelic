package gorelic

import (
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

func addMemoryMericsToComponent(component newrelic_platform_go.IComponent, ds DataSource, pollInterval int) {
	metrics.RegisterRuntimeMemStats(ds)
	metrics.CaptureRuntimeMemStatsOnce(ds)
	go metrics.CaptureRuntimeMemStats(ds, time.Duration(pollInterval)*time.Second)

	basePath := "Runtime/Memory/"
	curPath := basePath + "InUse/"
	units := "bytes"

	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.Alloc", curPath, "Total", units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapAlloc", curPath, "Heap", units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.StackInuse", curPath, "Stack", units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.MSpanInuse", curPath, "MSpanInuse", units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.MCacheInuse", curPath, "MCacheInuse", units))

	curPath = basePath + "Operations/"
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Lookups", curPath, "NoPointerLookups", "lookups"))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Mallocs", curPath, "NoMallocs", "mallocs"))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Frees", curPath, "NoFrees", "frees"))

	curPath = basePath + "SysMem/"
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Sys", curPath, "Total", units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.HeapSys", curPath, "Heap", units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.StackSys", curPath, "Stack", units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.MSpanSys", curPath, "Mspan", units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.MCacheSys", curPath, "MCache", units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.BuckHashSys", curPath, "BuckHash", units))
}
