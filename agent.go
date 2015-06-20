package gorelic

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/courtf/go-metrics"
	"github.com/courtf/newrelic_platform_go"
)

const (
	// DefaultNewRelicPollInterval - how often we will report metrics to NewRelic.
	// Recommended values is 60 seconds
	DefaultNewRelicPollInterval = 60

	// DefaultGcPollIntervalInSeconds - how often we will get garbage collector run statistic
	// Default value is - every 10 seconds
	// During GC stat pooling - mheap will be locked, so be carefull changing this value
	DefaultGcPollIntervalInSeconds = 10

	// DefaultMemoryAllocatorPollIntervalInSeconds - how often we will get memory allocator statistic.
	// Default value is - every 60 seconds
	// During this process stoptheword() is called, so be carefull changing this value
	DefaultMemoryAllocatorPollIntervalInSeconds = 60

	//DefaultAgentGuid is plugin ID in NewRelic.
	//You should not change it unless you want to create your own plugin.
	DefaultAgentGuid = "com.acmeaom.GoPlugin"

	//CurrentAgentVersion is plugin version
	CurrentAgentVersion = "0.0.1"

	//DefaultAgentName in NewRelic GUI. You can change it.
	DefaultAgentName = "Go Plugin"

	httpThroughPutDataSourceKey = "gorelic.http.throughput"
	httpStatusDataSourceKey     = "gorelic.http.status." // add code to the end
)

//Agent - is NewRelic agent implementation.
//Agent start separate go routine which will report data to NewRelic
type Agent struct {
	NewrelicName                string
	NewrelicLicense             string
	NewrelicPollInterval        int
	Verbose                     bool
	CollectGcStat               bool
	CollectMemoryStat           bool
	CollectHTTPStat             bool
	CollectHTTPStatuses         bool
	GCPollInterval              int
	MemoryAllocatorPollInterval int
	AgentGUID                   string
	AgentVersion                string
	plugin                      *newrelic_platform_go.NewrelicPlugin
	HTTPTimer                   metrics.Timer
	Tracer                      *Tracer
	CustomMetrics               []newrelic_platform_go.IMetrica
	cmLk                        sync.Mutex
	running                     uint32

	// All HTTP requests will be done using this client. Change it if you need
	// to use a proxy.
	Client http.Client

	// data source for internal use
	dataSource DataSource
}

// NewAgent builds new Agent objects.
func NewAgent() *Agent {
	agent := &Agent{
		NewrelicName:                DefaultAgentName,
		NewrelicPollInterval:        DefaultNewRelicPollInterval,
		Verbose:                     false,
		CollectGcStat:               true,
		CollectMemoryStat:           true,
		GCPollInterval:              DefaultGcPollIntervalInSeconds,
		MemoryAllocatorPollInterval: DefaultMemoryAllocatorPollIntervalInSeconds,
		AgentGUID:                   DefaultAgentGuid,
		AgentVersion:                CurrentAgentVersion,
		Tracer:                      nil,
		CustomMetrics:               make([]newrelic_platform_go.IMetrica, 0),
		dataSource:                  NewDataSource(metrics.NewRegistry()),
	}
	return agent
}

// used by proxyWrapper below to record http statuses
type statusRecorder struct {
	http.ResponseWriter
	writeHeader func(status int)
}

func (sr statusRecorder) WriteHeader(status int) {
	sr.writeHeader(status)
}

type proxyWrapper struct {
	*tHTTPHandler
	serveHTTP func(w http.ResponseWriter, req *http.Request)
}

func (pw proxyWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pw.serveHTTP(w, req)
}

func newProxyWrapper(agent *Agent, proxy *tHTTPHandler) proxyWrapper {
	return proxyWrapper{
		proxy,
		func(w http.ResponseWriter, req *http.Request) {
			proxy.ServeHTTP(
				statusRecorder{
					w,
					func(status int) {
						go agent.dataSource.IncCounterForKey(statusKeyFunc(status), 1)
						w.WriteHeader(status)
					},
				},
				req,
			)
		},
	}
}

//WrapHTTPHandlerFunc  instrument HTTP handler functions to collect HTTP metrics
func (agent *Agent) WrapHTTPHandlerFunc(h tHTTPHandlerFunc) tHTTPHandlerFunc {
	agent.CollectHTTPStat = true
	agent.initTimer()
	proxy := newHTTPHandlerFunc(h)
	proxy.timer = agent.HTTPTimer

	if agent.CollectHTTPStatuses {
		pr := newProxyWrapper(agent, proxy)
		return pr.ServeHTTP
	}

	return proxy.ServeHTTP
}

//WrapHTTPHandler  instrument HTTP handler object to collect HTTP metrics
func (agent *Agent) WrapHTTPHandler(h http.Handler) http.Handler {
	agent.CollectHTTPStat = true
	agent.initTimer()

	proxy := newHTTPHandler(h)
	proxy.timer = agent.HTTPTimer

	if agent.CollectHTTPStatuses {
		return newProxyWrapper(agent, proxy)
	}

	return proxy
}

//AddCustomMetric adds metric to be collected periodically with NewrelicPollInterval interval
func (agent *Agent) AddCustomMetric(metric newrelic_platform_go.IMetrica) {
	agent.cmLk.Lock()
	agent.CustomMetrics = append(agent.CustomMetrics, metric)
	agent.cmLk.Unlock()

	if atomic.LoadUint32(&agent.running) > 0 {
		// the plugin is only modified in agent.Run, where a single component is added
		agent.plugin.ComponentModels[0].AddMetrica(metric)
	}
}

//Run initialize Agent instance and start harvest go routine
func (agent *Agent) Run() error {
	if agent.NewrelicLicense == "" {
		return errors.New("please, pass a valid newrelic license key")
	}

	agent.plugin = newrelic_platform_go.NewNewrelicPlugin(agent.AgentVersion, agent.NewrelicLicense, agent.NewrelicPollInterval)
	agent.plugin.Client = agent.Client

	var component newrelic_platform_go.IComponent
	component = newrelic_platform_go.NewPluginComponent(agent.NewrelicName, agent.AgentGUID)

	// Add default metrics and tracer.
	addRuntimeMetricsToComponent(component)
	agent.Tracer = newTracer(component, agent.dataSource)

	// Check agent flags and add relevant metrics.
	if agent.CollectGcStat {
		addGCMericsToComponent(component, agent.dataSource, agent.GCPollInterval)
		agent.debug(fmt.Sprintf("Init GC metrics collection. Poll interval %d seconds.", agent.GCPollInterval))
	}

	if agent.CollectMemoryStat {
		addMemoryMericsToComponent(component, agent.dataSource, agent.MemoryAllocatorPollInterval)
		agent.debug(fmt.Sprintf("Init memory allocator metrics collection. Poll interval %d seconds.", agent.MemoryAllocatorPollInterval))
	}

	if agent.CollectHTTPStat {
		agent.initTimer()
		addHTTPMericsToComponent(component, agent.dataSource, httpThroughPutDataSourceKey)
		agent.debug(fmt.Sprintf("Init HTTP metrics collection."))
	}

	if agent.CollectHTTPStatuses {
		statuses := getHTTPStatuses()
		agent.initStatusCounters(statuses)
		addHTTPStatusMetricsToComponent(component, agent.dataSource, statuses, statusKeyFunc)
		agent.debug(fmt.Sprintf("Init HTTP status metrics collection."))
	}

	// Init newrelic reporting plugin.
	agent.plugin = newrelic_platform_go.NewNewrelicPlugin(agent.AgentVersion, agent.NewrelicLicense, agent.NewrelicPollInterval)
	agent.plugin.Verbose = agent.Verbose

	agent.cmLk.Lock()
	for _, metric := range agent.CustomMetrics {
		component.AddMetrica(metric)
		agent.debug(fmt.Sprintf("Init %s metric collection.", metric.GetName()))
	}

	// Add our metrics component to the plugin.
	agent.plugin.AddComponent(component)

	atomic.StoreUint32(&agent.running, 1)
	agent.cmLk.Unlock()

	// Start reporting!
	agent.plugin.Run()
	return nil
}

//Initialize global metrics.Timer object, used to collect HTTP metrics
func (agent *Agent) initTimer() {
	if agent.HTTPTimer == nil {
		agent.HTTPTimer = metrics.NewTimer()
		agent.dataSource.Register(httpThroughPutDataSourceKey, agent.HTTPTimer)
	}
}

//Initialize metrics.Counters objects, used to collect HTTP statuses
func (agent *Agent) initStatusCounters(statuses []int) {
	for _, statusCode := range statuses {
		agent.dataSource.Register(statusKeyFunc(statusCode), metrics.NewCounter())
	}
}

func getHTTPStatuses() []int {
	return []int{
		http.StatusContinue, http.StatusSwitchingProtocols,

		http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNonAuthoritativeInfo,
		http.StatusNoContent, http.StatusResetContent, http.StatusPartialContent,

		http.StatusMultipleChoices, http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther,
		http.StatusNotModified, http.StatusUseProxy, http.StatusTemporaryRedirect,

		http.StatusBadRequest, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusForbidden,
		http.StatusNotFound, http.StatusMethodNotAllowed, http.StatusNotAcceptable, http.StatusProxyAuthRequired,
		http.StatusRequestTimeout, http.StatusConflict, http.StatusGone, http.StatusLengthRequired,
		http.StatusPreconditionFailed, http.StatusRequestEntityTooLarge, http.StatusRequestURITooLong, http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable, http.StatusExpectationFailed, http.StatusTeapot,

		http.StatusInternalServerError, http.StatusNotImplemented, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusHTTPVersionNotSupported,
	}
}

func statusKeyFunc(status int) string {
	return httpStatusDataSourceKey + fmt.Sprintf("%d", status)
}

//Print debug messages
func (agent *Agent) debug(msg string) {
	if agent.Verbose {
		log.Println(msg)
	}
}
