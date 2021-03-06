package gorelic

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/courtf/newrelic_platform_go"
)

const linuxSystemQueryInterval = 60

// Number of goroutines metrica
type noGoroutinesMetrica struct{}

func (metrica *noGoroutinesMetrica) GetName() string {
	return "Runtime/General/NOGoroutines"
}
func (metrica *noGoroutinesMetrica) GetUnits() string {
	return "goroutines"
}
func (metrica *noGoroutinesMetrica) GetValue() (float64, error) {
	return float64(runtime.NumGoroutine()), nil
}
func (metrica *noGoroutinesMetrica) ClearSentData() {
	// no-op
}

// Number of CGO calls metrica
type noCgoCallsMetrica struct {
	lastValue int64
}

func (metrica *noCgoCallsMetrica) GetName() string {
	return "Runtime/General/NOCgoCalls"
}
func (metrica *noCgoCallsMetrica) GetUnits() string {
	return "calls"
}
func (metrica *noCgoCallsMetrica) GetValue() (float64, error) {
	currentValue := runtime.NumCgoCall()
	value := float64(currentValue - metrica.lastValue)
	metrica.lastValue = currentValue

	return value, nil
}
func (metrica *noCgoCallsMetrica) ClearSentData() {
	// no-op
}

//OS specific metrics data source interface
type iSystemDataSource interface {
	GetValue(key string) (float64, error)
}

// iSystemDataSource fabrica
func newSystemDataSource() iSystemDataSource {
	var ds iSystemDataSource
	switch runtime.GOOS {
	default:
		ds = &systemDataSource{}
	case "linux":
		ds = &linuxSystemDataSource{
			systemData: make(map[string]string),
		}
	}
	return ds
}

//Default implementation of iSystemDataSource. Just return an error
type systemDataSource struct{}

func (ds *systemDataSource) GetValue(key string) (float64, error) {
	return 0, fmt.Errorf("this metrica was not implemented yet for %s", runtime.GOOS)
}

// Linux OS implementation of ISystemDataSource
type linuxSystemDataSource struct {
	lastUpdate time.Time
	systemData map[string]string
}

func (ds *linuxSystemDataSource) GetValue(key string) (float64, error) {
	if err := ds.checkAndUpdateData(); err != nil {
		return 0, err
	} else if val, ok := ds.systemData[key]; !ok {
		return 0, fmt.Errorf("system data with key %s was not found", key)
	} else if key == "VmSize" || key == "VmPeak" || key == "VmHWM" || key == "VmRSS" {
		valueParts := strings.Split(val, " ")
		if len(valueParts) != 2 {
			return 0, fmt.Errorf("invalid format for value %s", key)
		}
		valConverted, err := strconv.ParseFloat(valueParts[0], 64)
		if err != nil {
			return 0, err
		}
		switch valueParts[1] {
		case "kB":
			valConverted *= 1 << 10
		case "mB":
			valConverted *= 1 << 20
		case "gB":
			valConverted *= 1 << 30
		}
		return valConverted, nil
	} else if valConverted, err := strconv.ParseFloat(val, 64); err != nil {
		return valConverted, nil
	} else {
		return valConverted, nil
	}
}
func (ds *linuxSystemDataSource) checkAndUpdateData() error {
	startTime := time.Now()
	if startTime.Sub(ds.lastUpdate) > time.Second*linuxSystemQueryInterval {
		path := fmt.Sprintf("/proc/%d/status", os.Getpid())
		rawStats, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(rawStats), "\n")
		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])

				ds.systemData[k] = v
			}
		}
		ds.lastUpdate = startTime
	}
	return nil
}

// OS specific metrica
type systemMetrica struct {
	sourceKey    string
	newrelicName string
	units        string
	dataSource   iSystemDataSource
}

func (metrica *systemMetrica) GetName() string {
	return metrica.newrelicName
}
func (metrica *systemMetrica) GetUnits() string {
	return metrica.units
}
func (metrica *systemMetrica) GetValue() (float64, error) {
	return metrica.dataSource.GetValue(metrica.sourceKey)
}
func (metrica *systemMetrica) ClearSentData() {
	// no-op
}

func addRuntimeMetricsToComponent(component newrelic_platform_go.IComponent) {
	component.AddMetrica(&noGoroutinesMetrica{})
	component.AddMetrica(&noCgoCallsMetrica{})

	ds := newSystemDataSource()
	metrics := []*systemMetrica{
		{
			sourceKey:    "Threads",
			units:        "Threads",
			newrelicName: "Runtime/System/Threads",
		},
		{
			sourceKey:    "FDSize",
			units:        "fd",
			newrelicName: "Runtime/System/FDSize",
		},
		// Peak virtual memory size
		{
			sourceKey:    "VmPeak",
			units:        "bytes",
			newrelicName: "Runtime/System/Memory/VmPeakSize",
		},
		//Virtual memory size
		{
			sourceKey:    "VmSize",
			units:        "bytes",
			newrelicName: "Runtime/System/Memory/VmCurrent",
		},
		//Peak resident set size
		{
			sourceKey:    "VmHWM",
			units:        "bytes",
			newrelicName: "Runtime/System/Memory/RssPeak",
		},
		//Resident set size
		{
			sourceKey:    "VmRSS",
			units:        "bytes",
			newrelicName: "Runtime/System/Memory/RssCurrent",
		},
	}
	for _, m := range metrics {
		m.dataSource = ds
		component.AddMetrica(m)
	}
}
