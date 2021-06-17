package collector

import (
	"context"
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"sync"
	"time"
)

// Metric name parts.
const (
	// Exporter namespace.
	namespace = "vsphere"
	// Subsystem(s).
	subsystem = "exporter"
	// Math constant for picoseconds to seconds.
	picoSeconds = 1e12
)

// Metric descriptors.
var (
	totalScrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "collector_duration_seconds"),
		"Collector time duration.",
		nil, nil,
	)
)

// Exporter collects redfish metrics. It implements prometheus.Collector.
type VshpereCollector struct {
	vsClient   *vmware.VMClient
	collectors map[string]prometheus.Collector
	vsherehUp  prometheus.Gauge
}

func NewVshpereCollector(context context.Context, url string, username string, password string) *VshpereCollector {
	var collectors map[string]prometheus.Collector

	vsClient := vmware.NewVMClient(context, url, username, password)
	if vsClient == nil {
		log.Infof("Errors occour when creating vshpere client")
	} else {
		hostCollector := NewHostCollector(namespace, vsClient)
		collectors = map[string]prometheus.Collector{"host": hostCollector}
	}

	return &VshpereCollector{
		vsClient:   vsClient,
		collectors: collectors,
		vsherehUp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "up",
				Help:      "vsphere up",
			},
		),
	}
}

// Describe implements prometheus.Collector.
func (r *VshpereCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range r.collectors {
		collector.Describe(ch)
	}

}

// Collect implements prometheus.Collector.
func (r *VshpereCollector) Collect(ch chan<- prometheus.Metric) {

	scrapeTime := time.Now()
	if r.vsClient != nil {
		defer r.vsClient.Logout()
		r.vsherehUp.Set(1)
		wg := &sync.WaitGroup{}
		wg.Add(len(r.collectors))

		defer wg.Wait()
		for _, collector := range r.collectors {
			go func(collector prometheus.Collector) {
				defer wg.Done()
				collector.Collect(ch)
			}(collector)
		}
	} else {
		r.vsherehUp.Set(0)
	}

	ch <- r.vsherehUp
	ch <- prometheus.MustNewConstMetric(totalScrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds())
}
