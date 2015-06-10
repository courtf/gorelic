package gorelic

import (
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

func addHTTPMericsToComponent(component newrelic_platform_go.IComponent, timer metrics.Timer) {
	addMeterMetrics(component, timer, "http/throughput", "rps")
	addTimedHistogramMetrics(component, timer, "http/throughput")
}
