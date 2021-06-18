package collector

import (
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	vmSubsystem  = "vm"
	vmMapping    = map[string]string{}
	vmLabelNames = []string{"name", "guest", "host"}
	//vmLabelNames = []string{"category"}
	vmMetrics = map[string]vmMetric{
		"vm_uptime": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, vmSubsystem, "uptime"),
				"the virtual machine uptime in seconds",
				vmLabelNames,
				nil,
			),
		},
	}
)

// A VmCollector implements the prometheus.Collector.
type VmCollector struct {
	vsClient              *vmware.VMClient
	metrics               map[string]vmMetric
	collectorScrapeStatus *prometheus.GaugeVec
}

type vmMetric struct {
	desc *prometheus.Desc
}

// NewVmCollector returns a collector that collecting vm statistics
func NewVmCollector(namespace string, vsClient *vmware.VMClient) *VmCollector {

	// get service from redfish client

	return &VmCollector{
		vsClient: vsClient,
		metrics:  vmMetrics,
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

func (v *VmCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range v.metrics {
		ch <- metric.desc
	}
	v.collectorScrapeStatus.Describe(ch)

}

func (v *VmCollector) Collect(ch chan<- prometheus.Metric) {
	// get a vm list from vsphere client
	if vmList, err := v.vsClient.ListVirtualMachine(); err != nil {
		log.Infof("Errors Getting vm list from vsphere : %s", err)
	} else {
		// process the vm status
		for _, vm := range vmList {
			vmSummary := vm.Summary
			vmConfig := vm.Config
			vmQuickStats := vmSummary.QuickStats
			vmName := vmConfig.Name
			vmID := vm.ManagedEntity.ExtensibleManagedObject.Self.Value
			vmMapping[vmID] = vmName
			vmGuestFullName := vmConfig.GuestFullName
			vmHost := hostMapping[vmSummary.Runtime.Host.Value]
			vmLabelValues := []string{vmName, vmGuestFullName, vmHost}
			//
			vmUptimeValue := float64(vmQuickStats.UptimeSeconds)

			ch <- prometheus.MustNewConstMetric(v.metrics["vm_uptime"].desc, prometheus.GaugeValue, vmUptimeValue, vmLabelValues...)

		}

		v.collectorScrapeStatus.WithLabelValues("virtualmachine").Set(float64(1))
	}
}
