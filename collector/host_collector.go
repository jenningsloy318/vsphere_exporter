package collector

import (
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vmware/govmomi/vim25/types"
)

var (
	hostSubsystem = "host"
	//hostLabelNames = []string{"category", "name"}
	hostLabelNames = []string{"category"}
	hostMetrics    = map[string]hostMetric{
		"host_powerstate": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "up"),
				"up state of the host, 1(OK),2(Down)",
				hostLabelNames,
				nil,
			),
		},
	}
)

// A HostCollector implements the prometheus.Collector.
type HostCollector struct {
	vsClient              *vmware.VMClient
	metrics               map[string]hostMetric
	collectorScrapeStatus *prometheus.GaugeVec
}

type hostMetric struct {
	desc *prometheus.Desc
}

// NewHostCollector returns a collector that collecting host statistics
func NewHostCollector(namespace string, vsClient *vmware.VMClient) *HostCollector {

	// get service from redfish client

	return &HostCollector{
		vsClient: vsClient,
		metrics:  hostMetrics,
		collectorScrapeStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "collector_scrape_status",
				Help:      "collector_scrape_status",
			},
			[]string{"collector"},
		),
	}
}

func (h *HostCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range h.metrics {
		ch <- metric.desc
	}
	h.collectorScrapeStatus.Describe(ch)

}

func (h *HostCollector) Collect(ch chan<- prometheus.Metric) {
	// get a host list from vsphere client
	if hostList, err := h.vsClient.ListHost(); err != nil {
		log.Infof("Errors Getting host list from vsphere : %s", err)
	} else {
		// process the host status
		for _, host := range hostList {
			//hostSummary := host.Summary
			//hostHW := host.Hardware
			hostRuntime := host.Runtime
			//hostConfig := host.Config
			//hostLabelValues := []string{"ESXi", string(hostConfig.Host.Value)}
			hostLabelValues := []string{"ESXi"}
			hostPowerStateValue := retrievePowerState(hostRuntime.PowerState)

			ch <- prometheus.MustNewConstMetric(h.metrics["host_powerstate"].desc, prometheus.GaugeValue, hostPowerStateValue, hostLabelValues...)

		}
		h.collectorScrapeStatus.WithLabelValues("host").Set(float64(1))
	}
}

func retrievePowerState(s types.HostSystemPowerState) float64 {

	if s == types.HostSystemPowerStatePoweredOn {
		return float64(1)
	}
	if s == types.HostSystemPowerStatePoweredOff {
		return float64(2)
	}
	if s == types.HostSystemPowerStatePoweredOff {
		return float64(3)
	}

	return float64(4)
}
