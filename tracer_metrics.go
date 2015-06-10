package gorelic

import (
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

type Tracer struct {
	metrics   map[string]*TraceTransaction
	component newrelic_platform_go.IComponent
}

func newTracer(component newrelic_platform_go.IComponent) *Tracer {
	return &Tracer{make(map[string]*TraceTransaction), component}
}

func (t *Tracer) Trace(name string, traceFunc func()) {
	trace := t.BeginTrace(name)
	defer trace.EndTrace()
	traceFunc()
}

func (t *Tracer) BeginTrace(name string) *Trace {
	tracerName := "Trace/" + name
	m := t.metrics[tracerName]
	if m == nil {
		t.metrics[tracerName] = &TraceTransaction{tracerName, metrics.NewTimer()}
		m = t.metrics[tracerName]
		m.addMetricsToComponent(t.component)
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
	name  string
	timer metrics.Timer
}

func (transaction *TraceTransaction) addMetricsToComponent(component newrelic_platform_go.IComponent) {
	addTimedHistogramMetrics(component, transaction.timer, transaction.name)
}
