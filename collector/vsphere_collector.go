package collector

import (
	"context"
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vmware/govmomi/vim25/types"
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

	vsClient, err := vmware.NewVMClient(context, url, username, password)
	if err != nil {
		log.Errorf("Errors occour when creating vshpere client, %v", err)
	} else {
		hostCollector := NewHostCollector(namespace, vsClient)
		vmCollector := NewVmCollector(namespace, vsClient)
		collectors = map[string]prometheus.Collector{"host": hostCollector, "vm": vmCollector}
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

		r.vsherehUp.Set(1)

		r.collectors["host"].Collect(ch)
		r.collectors["vm"].Collect(ch)
	} else {
		r.vsherehUp.Set(0)
	}

	ch <- r.vsherehUp
	ch <- prometheus.MustNewConstMetric(totalScrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds())
	r.vsClient.Logout()
}

func parseOveralStatus(status types.ManagedEntityStatus) float64 {
	if status == "green" {
		return float64(1)
	}
	if status == "yellow" {
		return float64(2)
	}
	if status == "gray" {
		return float64(3)
	}

	return float64(4)
}

func parsePowerState(powerState types.HostSystemPowerState) float64 {
	if powerState == "poweredOn" {
		return float64(1)
	}
	if powerState == "poweredOff" {
		return float64(2)
	}
	if powerState == "standBy" {
		return float64(3)
	}

	return float64(4)
}
func parseConnectionState(connectionState types.HostSystemConnectionState) float64 {
	if connectionState == "active" {
		return float64(1)
	}
	if connectionState == "activeDefer" {
		return float64(2)
	}
	if connectionState == "armed" {
		return float64(3)
	}
	if connectionState == "init" {
		return float64(4)
	}
	if connectionState == "down" {
		return float64(5)
	}
	return float64(6)
}

func parseHostStandbyMode(standbyMode string) float64 {
	if standbyMode == "in" {
		return float64(1)
	}
	if standbyMode == "exiting" {
		return float64(2)
	}
	if standbyMode == "entering" {
		return float64(3)
	}
	return float64(4)
}
