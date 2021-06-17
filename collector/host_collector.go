package collector

import (
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	hostSubsystem  = "host"
	hostLabelNames = []string{"category", "name"}
	//hostLabelNames = []string{"category"}
	hostMetrics = map[string]hostMetric{
		"host_uptime": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "uptime"),
				"uptime of the",
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
			hostSummary := host.Summary
			hostName := hostSummary.Config.Name
			hostLabelValues := []string{"ESXi", hostName}
			hostUptimeValue := float64(hostSummary.QuickStats.Uptime)

			ch <- prometheus.MustNewConstMetric(h.metrics["host_uptime"].desc, prometheus.GaugeValue, hostUptimeValue, hostLabelValues...)

		}
		h.collectorScrapeStatus.WithLabelValues("host").Set(float64(1))
	}
}
