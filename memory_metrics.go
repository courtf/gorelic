package gorelic

import (
	"path/filepath"
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

func addMemoryMetricsToComponent(component newrelic_platform_go.IComponent, ds DataSource, pollInterval int) {
	metrics.RegisterRuntimeMemStats(ds)
	metrics.CaptureRuntimeMemStatsOnce(ds)
	go metrics.CaptureRuntimeMemStats(ds, time.Duration(pollInterval)*time.Second)

	basePath := "Runtime/Memory/"
	curPath := basePath + "InUse/"
	units := "bytes"

	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.TotalAlloc", filepath.Join(curPath, "TotalAlloc"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.TotalAllocSince", filepath.Join(curPath, "TotalAllocDelta"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.Alloc", filepath.Join(curPath, "Alloc"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.AllocSince", filepath.Join(curPath, "AllocDelta"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapInuse", filepath.Join(curPath, "HeapInuse"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.StackInuse", filepath.Join(curPath, "Stack"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.MSpanInuse", filepath.Join(curPath, "MSpanInuse"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.MCacheInuse", filepath.Join(curPath, "MCacheInuse"), units))

	curPath = basePath + "Heap/"
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.NumGC", filepath.Join(curPath, "NumGC"), "calls"))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapSys", filepath.Join(curPath, "HeapSys"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapAlloc", filepath.Join(curPath, "HeapAlloc"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapIdle", filepath.Join(curPath, "HeapIdle"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapInuse", filepath.Join(curPath, "HeapInuse"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapObjects", filepath.Join(curPath, "HeapObjects"), units))
	component.AddMetrica(NewGaugeMetrica(ds, "runtime.MemStats.HeapReleased", filepath.Join(curPath, "HeapReleased"), units))

	curPath = basePath + "Operations/"
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Lookups", filepath.Join(curPath, "NoPointerLookups"), "lookups"))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Mallocs", filepath.Join(curPath, "NoMallocs"), "mallocs"))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Frees", filepath.Join(curPath, "NoFrees"), "frees"))

	curPath = basePath + "SysMem/"
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.Sys", filepath.Join(curPath, "Total"), units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.HeapSys", filepath.Join(curPath, "Heap"), units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.StackSys", filepath.Join(curPath, "Stack"), units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.MSpanSys", filepath.Join(curPath, "Mspan"), units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.MCacheSys", filepath.Join(curPath, "MCache"), units))
	component.AddMetrica(NewGaugeDeltaMetrica(ds, "runtime.MemStats.BuckHashSys", filepath.Join(curPath, "BuckHash"), units))
}
