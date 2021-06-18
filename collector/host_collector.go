package collector

import (
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	hostSubsystem  = "host"
	hostMapping    = map[string]string{}
	hostLabelNames = []string{"hostname", "os"}
	//hostLabelNames = []string{"category"}
	hostMetrics = map[string]hostMetric{

		"host_connection_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "connection_state"),
				"host connection state to vcenter, 1 for active, 2 for activeDefer, 3 for armed, 4 for init,5 for down,6 for unkown",
				hostLabelNames,
				nil,
			),
		},
		"host_power_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "power_state"),
				"host power state, 1 for poweron, 2 for poweroff, 3 for standby, 4 for unkown",
				hostLabelNames,
				nil,
			),
		},

		"host_standby_mode": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "standby_mode"),
				"host standby mode, 1 for in , 2 for exiting, 3 for entering, 4 for none",
				hostLabelNames,
				nil,
			),
		},
		"host_maintenance_mode": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "maintenance_mode"),
				"if the host is in maintenance mode",
				hostLabelNames,
				nil,
			),
		},
		"host_in_quarantine_mode": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "in_quarantine_mode"),
				"if the host is in quarantine mode",
				hostLabelNames,
				nil,
			),
		},
		"host_vmotion_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "vmotion_status"),
				"the status of vmotion, 1 is enabled, 0 is disabled",
				hostLabelNames,
				nil,
			),
		},
		"host_fault_tolerance_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "fault_tolerance_status"),
				"the status of vmotion, 1 is enabled, 0 is disabled",
				hostLabelNames,
				nil,
			),
		},
		"host_uptime": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "uptime"),
				"uptime of the host",
				hostLabelNames,
				nil,
			),
		},
		"host_overall_cpu_used": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "overall_cpu_used"),
				"host overall cpu used in mhz",
				hostLabelNames,
				nil,
			),
		},
		"host_overall_memory_used": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "overall_memory_used"),
				"host overall memory used in MB",
				hostLabelNames,
				nil,
			),
		},
		"host_distributed_cpu_fairness": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "distributed_cpu_fairness"),
				"host distributed cpu fairness",
				hostLabelNames,
				nil,
			),
		},
		"host_distributed_memory_fairness": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "distributed_memory_fairness"),
				"host distributed memory fairness",
				hostLabelNames,
				nil,
			),
		},
		"host_available_pmem_capacity": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "available_pmem_capacity"),
				"host available pmem capacity",
				hostLabelNames,
				nil,
			),
		},
		"host_overall_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "overall_status"),
				"host overall status, 1 for green, 2 for yellow, 3 for gray, 4 for red",
				hostLabelNames,
				nil,
			),
		},
		"host_memory_size": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "memory_size"),
				"host memory size",
				hostLabelNames,
				nil,
			),
		},
		"host_cpu_sockets": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "cpu_sockets"),
				"host cpu socket number",
				hostLabelNames,
				nil,
			),
		},
		"host_cpu_cores": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "cpu_cores"),
				"host cpu cores",
				hostLabelNames,
				nil,
			),
		},
		"host_cpu_threads": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "cpu_threads"),
				"host cpu threads",
				hostLabelNames,
				nil,
			),
		},
		"host_nic_counts": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "nic_counts"),
				"host nic counts",
				hostLabelNames,
				nil,
			),
		},
		"host_hba_counts": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, hostSubsystem, "hba_counts"),
				"host hba counts",
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
			hostRumtime := host.Runtime
			hostName := hostSummary.Config.Name
			hostID := host.ManagedEntity.ExtensibleManagedObject.Self.Value
			hostMapping[hostID] = hostName
			esxiFullName := hostSummary.Config.Product.FullName
			hostLabelValues := []string{hostName, esxiFullName}

			// retrieve the connection state between host and vcenter
			hostConnectionStateValue := parseConnectionState(hostRumtime.ConnectionState)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_connection_state"].desc, prometheus.GaugeValue, hostConnectionStateValue, hostLabelValues...)

			// retrueve the powerstate
			hostPowerStateValue := parsePowerState(hostRumtime.PowerState)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_power_state"].desc, prometheus.GaugeValue, hostPowerStateValue, hostLabelValues...)
			// retrueve standby mode
			hostStandbyModeValue := parseHostStandbyMode(hostRumtime.StandbyMode)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_standby_mode"].desc, prometheus.GaugeValue, hostStandbyModeValue, hostLabelValues...)

			// retrieve the maintenance mode
			var hostMaintenanceModeValue float64
			if hostRumtime.InMaintenanceMode {
				hostMaintenanceModeValue = float64(1)
			} else {
				hostMaintenanceModeValue = float64(0)
			}
			ch <- prometheus.MustNewConstMetric(h.metrics["host_maintenance_mode"].desc, prometheus.GaugeValue, hostMaintenanceModeValue, hostLabelValues...)

			// retrieve the in quarantine  mode
			var hostInQuarantineModeValue float64
			if *hostRumtime.InQuarantineMode {
				hostInQuarantineModeValue = float64(1)
			} else {
				hostInQuarantineModeValue = float64(0)
			}
			ch <- prometheus.MustNewConstMetric(h.metrics["host_in_quarantine_mode"].desc, prometheus.GaugeValue, hostInQuarantineModeValue, hostLabelValues...)

			// retrieve the vmotion status
			var hostVmotionStatusValue float64
			if hostSummary.Config.VmotionEnabled {
				hostVmotionStatusValue = float64(1)
			} else {
				hostVmotionStatusValue = float64(0)
			}

			ch <- prometheus.MustNewConstMetric(h.metrics["host_vmotion_status"].desc, prometheus.GaugeValue, hostVmotionStatusValue, hostLabelValues...)

			// retrieve host fault tolerance status
			var hostFaultToleranceStatusValue float64
			if *hostSummary.Config.FaultToleranceEnabled {
				hostFaultToleranceStatusValue = float64(1)
			}
			hostFaultToleranceStatusValue = float64(0)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_fault_tolerance_status"].desc, prometheus.GaugeValue, hostFaultToleranceStatusValue, hostLabelValues...)

			hostQuickStats := hostSummary.QuickStats
			// retrieve the uptime of host
			hostUptimeValue := float64(hostQuickStats.Uptime)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_uptime"].desc, prometheus.GaugeValue, hostUptimeValue, hostLabelValues...)
			// retrieve the overall CPU usage in mhz
			hostOverallCpuUsedValue := float64(hostQuickStats.OverallCpuUsage)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_overall_cpu_used"].desc, prometheus.GaugeValue, hostOverallCpuUsedValue, hostLabelValues...)

			// retrieve the overall memory usage in mhz
			hostOverallMemoryUsedValue := float64(hostQuickStats.OverallMemoryUsage)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_overall_memory_used"].desc, prometheus.GaugeValue, hostOverallMemoryUsedValue, hostLabelValues...)

			// retrieve the distributed CPU fairness
			hostDistributedCpuFairnessValue := float64(hostQuickStats.DistributedCpuFairness)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_distributed_cpu_fairness"].desc, prometheus.GaugeValue, hostDistributedCpuFairnessValue, hostLabelValues...)

			// retrieve the distributed Memory fairness
			hostDistributedMemoryFairnessValue := float64(hostQuickStats.DistributedMemoryFairness)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_distributed_memory_fairness"].desc, prometheus.GaugeValue, hostDistributedMemoryFairnessValue, hostLabelValues...)

			// retrieve the available PMem capacity
			hostAvailablePMemCapacityValue := float64(hostQuickStats.AvailablePMemCapacity)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_available_pmem_capacity"].desc, prometheus.GaugeValue, hostAvailablePMemCapacityValue, hostLabelValues...)

			// retrueve the overall status

			hostOveralStatusValue := parseOveralStatus(hostSummary.OverallStatus)

			ch <- prometheus.MustNewConstMetric(h.metrics["host_overall_status"].desc, prometheus.GaugeValue, hostOveralStatusValue, hostLabelValues...)

			// retrueve the memory size
			hostHardware := hostSummary.Hardware
			hostMemSizeValue := float64(hostHardware.MemorySize)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_memory_size"].desc, prometheus.GaugeValue, hostMemSizeValue, hostLabelValues...)

			// retrueve the cpu counts
			hostCpuCountsValue := float64(hostHardware.NumCpuPkgs)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_cpu_sockets"].desc, prometheus.GaugeValue, hostCpuCountsValue, hostLabelValues...)

			// retrueve the cpu cores
			hostCpuCoresValue := float64(hostHardware.NumCpuCores)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_cpu_cores"].desc, prometheus.GaugeValue, hostCpuCoresValue, hostLabelValues...)

			// retrueve the cpu threads
			hostCputhreadsValue := float64(hostHardware.NumCpuThreads)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_cpu_threads"].desc, prometheus.GaugeValue, hostCputhreadsValue, hostLabelValues...)

			// retrueve the nic counts
			hostNicCountsValue := float64(hostHardware.NumNics)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_nic_counts"].desc, prometheus.GaugeValue, hostNicCountsValue, hostLabelValues...)

			// retrueve the hba counts
			hostHbaCountsValue := float64(hostHardware.NumHBAs)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_hba_counts"].desc, prometheus.GaugeValue, hostHbaCountsValue, hostLabelValues...)

		}

		h.collectorScrapeStatus.WithLabelValues("host").Set(float64(1))
	}
}
