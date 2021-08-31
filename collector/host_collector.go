package collector

import (
	"fmt"
	"github.com/jenningsloy318/vsphere_exporter/vmware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"strings"
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

			healthSystemRuntime := hostRumtime.HealthSystemRuntime

			systemHealthInfo := healthSystemRuntime.SystemHealthInfo
			hardwareStatusInfo := healthSystemRuntime.HardwareStatusInfo

			for _, hostNumericSensorInfo := range systemHealthInfo.NumericSensorInfo {
				sensorName := hostNumericSensorInfo.Name
				sensorID := hostNumericSensorInfo.Id
				sensorType := hostNumericSensorInfo.SensorType
				sensorCurrentValue := float64(hostNumericSensorInfo.CurrentReading)
				sensorDescriptionKey := hostNumericSensorInfo.HealthState.GetElementDescription().Key
				sensorDescription := hostNumericSensorInfo.HealthState.GetElementDescription().Description
				sensorDescriptionSummary := sensorDescription.Summary
				sensorDescriptionLabel := sensorDescription.Label
				sensorTimeStamp := hostNumericSensorInfo.TimeStamp
				fmt.Printf("sensorDescriptionKey:%s,sensorDescriptionSummary:%s,sensorDescriptionLabel:%s\n", sensorDescriptionKey, sensorDescriptionSummary, sensorDescriptionLabel)
				metricLabelNames := append(hostLabelNames, "sensorName", "sensorID", "sensorType", "Time")
				metricLabelValues := append(hostLabelValues, sensorName, sensorID, sensorType, sensorTimeStamp)
				metricName := fmt.Sprintf("%s_sensor_state", strings.ToLower(strings.ReplaceAll(sensorType, " ", "_")))

				sensorDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, metricName),
					sensorDescriptionSummary,
					metricLabelNames,
					nil,
				)
				ch <- prometheus.MustNewConstMetric(sensorDesc, prometheus.GaugeValue, sensorCurrentValue, metricLabelValues...)

			}

			memoryStatusInfo := hardwareStatusInfo.MemoryStatusInfo

			for _, memoryStatusInfoItem := range memoryStatusInfo {
				memoryStatusInfoData := memoryStatusInfoItem.GetHostHardwareElementInfo()
				memoryStatusName := memoryStatusInfoData.Name
				memoryStatusDescription := memoryStatusInfoData.Status.GetElementDescription()
				metricLabelNames := append(hostLabelNames, "component")
				metricLabelValues := append(hostLabelValues, memoryStatusName)
				memoryHWDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, "memory_hardware_status"),
					"memory hardware status, 0 for green, 1 for red",
					metricLabelNames,
					nil,
				)
				var memoryHWstateValue float64
				if memoryStatusDescription.Key == "Green" {
					memoryHWstateValue = float64(0)
				} else {
					memoryHWstateValue = float64(1)
				}
				ch <- prometheus.MustNewConstMetric(memoryHWDesc, prometheus.GaugeValue, memoryHWstateValue, metricLabelValues...)

			}
			cpuStatusInfo := hardwareStatusInfo.CpuStatusInfo
			for _, cpuStatusInfoItem := range cpuStatusInfo {
				cpuStatusInfoData := cpuStatusInfoItem.GetHostHardwareElementInfo()
				cpuStatusName := cpuStatusInfoData.Name
				cpuStatusDescription := cpuStatusInfoData.Status.GetElementDescription()
				metricLabelNames := append(hostLabelNames, "component")
				metricLabelValues := append(hostLabelValues, cpuStatusName)
				cpuHWDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, "cpu_hardware_status"),
					"cpu hardware status, 0 for green, 1 for red",
					metricLabelNames,
					nil,
				)
				var cpuHWstateValue float64
				if cpuStatusDescription.Key == "Green" {
					cpuHWstateValue = float64(0)
				} else {
					cpuHWstateValue = float64(1)
				}
				ch <- prometheus.MustNewConstMetric(cpuHWDesc, prometheus.GaugeValue, cpuHWstateValue, metricLabelValues...)

			}

			StorageStatusInfo := hardwareStatusInfo.StorageStatusInfo
			for _, storageStatusInfoItem := range StorageStatusInfo {
				storageStatusName := storageStatusInfoItem.HostHardwareElementInfo.Name
				storageStatusDescription := storageStatusInfoItem.HostHardwareElementInfo.Status.GetElementDescription()

				metricLabelNames := append(hostLabelNames, "component")
				metricLabelValues := append(hostLabelValues, storageStatusName)
				storageHWDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, "storage_hardware_status"),
					"storage hardware status, 0 for green, 1 for red",
					metricLabelNames,
					nil,
				)
				var storageHWstateValue float64
				if storageStatusDescription.Key == "Green" {
					storageHWstateValue = float64(0)
				} else {
					storageHWstateValue = float64(1)
				}
				ch <- prometheus.MustNewConstMetric(storageHWDesc, prometheus.GaugeValue, storageHWstateValue, metricLabelValues...)

			}

			networkRuntimeInfo := hostRumtime.NetworkRuntimeInfo

			netStackInstanceRuntimeInfo := networkRuntimeInfo.NetStackInstanceRuntimeInfo
			for _, netStackInstanceRuntimeInfoItem := range netStackInstanceRuntimeInfo {

				netStackInstanceKey := netStackInstanceRuntimeInfoItem.NetStackInstanceKey
				netStackInstanceState := netStackInstanceRuntimeInfoItem.State

				metricLabelNames := append(hostLabelNames, "component")
				metricLabelValues := append(hostLabelValues, netStackInstanceKey)
				netStackStateDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, "network_statck_state"),
					"network stack state, 1 for active, 0 for inactive",
					metricLabelNames,
					nil,
				)
				var netStackstateValue float64
				if netStackInstanceState == "active" {
					netStackstateValue = float64(1)
				} else {
					netStackstateValue = float64(0)
				}
				ch <- prometheus.MustNewConstMetric(netStackStateDesc, prometheus.GaugeValue, netStackstateValue, metricLabelValues...)

			}

			networkResourceRuntime := networkRuntimeInfo.NetworkResourceRuntime.PnicResourceInfo
			for _, pnicResourceInfoItem := range networkResourceRuntime {
				pnicDevice := pnicResourceInfoItem.PnicDevice
				pnicAvailableBandwidthForVMTraffic := pnicResourceInfoItem.AvailableBandwidthForVMTraffic
				pnicUnusedBandwidthForVMTraffic := pnicResourceInfoItem.UnusedBandwidthForVMTraffic

				metricLabelNames := append(hostLabelNames, "component")
				metricLabelValues := append(hostLabelValues, pnicDevice)
				pnicAvailableBandwidthForVMTrafficDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, "pnic_available_bandwidth_for_vm_traffic"),
					"pnic_available_bandwidth_for_vm_traffic",
					metricLabelNames,
					nil,
				)

				pnicUnusedBandwidthForVMTrafficDesc := prometheus.NewDesc(
					prometheus.BuildFQName(namespace, hostSubsystem, "pnic_unused_bandwidth_for_vm_traffic"),
					"pnic_unused_bandwidth_for_vm_traffic",
					metricLabelNames,
					nil,
				)

				ch <- prometheus.MustNewConstMetric(pnicAvailableBandwidthForVMTrafficDesc, prometheus.GaugeValue, float64(pnicAvailableBandwidthForVMTraffic), metricLabelValues...)
				ch <- prometheus.MustNewConstMetric(pnicUnusedBandwidthForVMTrafficDesc, prometheus.GaugeValue, float64(pnicUnusedBandwidthForVMTraffic), metricLabelValues...)

			}
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

			// retrieve the overall status

			hostOveralStatusValue := parseOveralStatus(hostSummary.OverallStatus)

			ch <- prometheus.MustNewConstMetric(h.metrics["host_overall_status"].desc, prometheus.GaugeValue, hostOveralStatusValue, hostLabelValues...)

			// retrieve the memory size
			hostHardware := hostSummary.Hardware
			hostMemSizeValue := float64(hostHardware.MemorySize)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_memory_size"].desc, prometheus.GaugeValue, hostMemSizeValue, hostLabelValues...)

			// retrieve the cpu counts
			hostCpuCountsValue := float64(hostHardware.NumCpuPkgs)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_cpu_sockets"].desc, prometheus.GaugeValue, hostCpuCountsValue, hostLabelValues...)

			// retrieve the cpu cores
			hostCpuCoresValue := float64(hostHardware.NumCpuCores)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_cpu_cores"].desc, prometheus.GaugeValue, hostCpuCoresValue, hostLabelValues...)

			// retrieve the cpu threads
			hostCputhreadsValue := float64(hostHardware.NumCpuThreads)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_cpu_threads"].desc, prometheus.GaugeValue, hostCputhreadsValue, hostLabelValues...)

			// retrieve the nic counts
			hostNicCountsValue := float64(hostHardware.NumNics)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_nic_counts"].desc, prometheus.GaugeValue, hostNicCountsValue, hostLabelValues...)

			// retrieve the hba counts
			hostHbaCountsValue := float64(hostHardware.NumHBAs)
			ch <- prometheus.MustNewConstMetric(h.metrics["host_hba_counts"].desc, prometheus.GaugeValue, hostHbaCountsValue, hostLabelValues...)

		}

		h.collectorScrapeStatus.WithLabelValues("host").Set(float64(1))
	}
}
