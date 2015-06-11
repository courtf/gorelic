package gorelic

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

type Tracer struct {
	metrics   map[string]*TraceTransaction
	component newrelic_platform_go.IComponent
	ds        DataSource
}

func newTracer(component newrelic_platform_go.IComponent, ds DataSource) *Tracer {
	return &Tracer{make(map[string]*TraceTransaction), component, ds}
}

func (t *Tracer) Trace(name string, traceFunc func()) {
	trace := t.BeginTrace(name)
	defer trace.EndTrace()
	traceFunc()
}

func (t *Tracer) BeginTrace(name string) *Trace {
	name = strings.Trim(name, "/")
	basePath := filepath.Join("Trace", name)

	m := t.metrics[basePath]
	if m == nil {
		srcKey := "gorelic.trace." + name
		timer := metrics.NewTimer()
		t.ds.Register(srcKey, timer)
		m = &TraceTransaction{timer, srcKey, basePath}
		t.metrics[basePath] = m
		m.addMetricsToComponent(t.component, t.ds)
	}
	return &Trace{m, time.Now()}
}

type Trace struct {
	transaction *TraceTransaction
	startTime   time.Time
}

func (t *Trace) EndTrace() {
	t.transaction.timer.UpdateSince(t.startTime)
}

type TraceTransaction struct {
	timer                   metrics.Timer
	dataSourceKey, basePath string
}

func (transaction *TraceTransaction) addMetricsToComponent(component newrelic_platform_go.IComponent, ds DataSource) {
	addTimerHistogramMetrics(component, ds, transaction.dataSourceKey, transaction.basePath)
}
