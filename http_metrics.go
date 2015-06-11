package gorelic

import (
	"fmt"
	"net/http"
	"time"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

type tHTTPHandlerFunc func(http.ResponseWriter, *http.Request)
type tHTTPHandler struct {
	originalHandler     http.Handler
	originalHandlerFunc tHTTPHandlerFunc
	isFunc              bool
	timer               metrics.Timer
}

var httpTimer metrics.Timer

func newHTTPHandlerFunc(h tHTTPHandlerFunc) *tHTTPHandler {
	return &tHTTPHandler{
		isFunc:              true,
		originalHandlerFunc: h,
	}
}

func newHTTPHandler(h http.Handler) *tHTTPHandler {
	return &tHTTPHandler{
		isFunc:          false,
		originalHandler: h,
	}
}

func (handler *tHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer handler.timer.UpdateSince(startTime)

	if handler.isFunc {
		handler.originalHandlerFunc(w, req)
	} else {
		handler.originalHandler.ServeHTTP(w, req)
	}
}

func addHTTPMericsToComponent(component newrelic_platform_go.IComponent, ds DataSource, timerKey string) {
	addTimerMeterMetrics(component, ds, timerKey, "HTTP/Throughput/", "rps")
	addTimerHistogramMetrics(component, ds, timerKey, "HTTP/Throughput/")
}

func addHTTPStatusMetricsToComponent(component newrelic_platform_go.IComponent, ds DataSource, statuses []int,
	keyFunc func(int) string) {
	for _, s := range statuses {
		component.AddMetrica(NewCounterMetrica(ds, keyFunc(s), "HTTP/Status/", fmt.Sprintf("%d", s), "count"))
	}
}
