# VMware vROps(vRealize Operations) Input Plugin

The VMware vROps plugin uses the vROps API to gather metrics.

- Projects
- VMs
- Tanzu Projects
- Tanzu VMs

## Configuration

```toml @sample.conf
# Collects performance metrics from vROps(vRealize Operations) services
[[inputs.vrops]]
  ## List of vROps URL to be monitored. These three lines must be uncommented
  ## and edited for the plugin to work.
  url = "https://vrealize.example.com"
  username = "vrops"
  password = "vrops"

  ## Available services are:
  ## "projects", "vms", "tanzuProjects", "tanzuVMs"
  enabled_services = ["projects", "vms", "tanzuProjects", "tanzuVMs"]

  ## Projects
  ## Typical Project Stat metrics (if omitted or empty, all metrics are collected)
  project_stat_key = []

  ## VMs
  ## Typical VMs Stat metrics (if omitted or empty, all metrics are collected)
  vm_stat_key = ["sys|poweredOn",
    "cpu|usage_average",
    "cpu|corecount_provisioned",
    "config|hardware|num_Cpu",
    "config|hardware|disk_Space",
    "mem|usage_average",
    "mem|guest_provisioned",
    "diskspace|provisionedSpace",
    "diskspace|used",
    "guestfilesystem|capacity_total",
    "guestfilesystem|freespace_total",
    "guestfilesystem|usage_total",
    "guestfilesystem|percentage_total",
    "net|transmitted_average",
    "net|received_average",
    "net|usage_average",
    "virtualDisk|read_average",
    "virtualDisk|write_average",
    "virtualDisk:Aggregate of all instances|totalLatency",
    "virtualDisk:Aggregate of all instances|totalReadLatency_average",
    "virtualDisk:Aggregate of all instances|totalWriteLatency_average",
    ]

  ## Tanzu Projects
  ## Typical Tanzu Project Stat metrics (if omitted or empty, all metrics are collected)
  tanzu_project_stat_key = []

  ## Tanzu VMs
  ## Typical Tanzu VMs Stat metrics (if omitted or empty, all metrics are collected)
  tanzu_vm_stat_key = ["sys|poweredOn",
    "cpu|usage_average",
    "cpu|corecount_provisioned",
    "config|hardware|num_Cpu",
    "config|hardware|disk_Space",
    "mem|usage_average",
    "mem|guest_provisioned",
    "diskspace|provisionedSpace",
    "diskspace|used",
    "guestfilesystem|capacity_total",
    "guestfilesystem|freespace_total",
    "guestfilesystem|usage_total",
    "guestfilesystem|percentage_total",
    "net|transmitted_average",
    "net|received_average",
    "net|usage_average",
    "virtualDisk|read_average",
    "virtualDisk|write_average",
    "virtualDisk|peak_vDisk_iops",
    "virtualDisk:Aggregate of all instances|totalLatency",
    "virtualDisk:Aggregate of all instances|totalReadLatency_average",
    "virtualDisk:Aggregate of all instances|totalWriteLatency_average",
    ]

  ## Tanzu Project - Resource Kind keys
  # tanzu_project_resource_kind = "ResourcePool"

  ## Tanzu VM - Resource Kind keys
  # tanzu_vm_resource_kind = "VirtualMachine"

  ## Tanzu Project Filter
  ## Object used to lookup Tanzu Project with various filtering criteria
  ## Indicates the conjunction of the filtering criteria
  ## Either all of the filtering criteria apply together (AND operation) or any of the filtering criteria could be applied (OR operation)
  ## Defaults to OR.
  # [inputs.vrops.tanzu_project_property_conditions]
  #   conjunctionOperator="AND"

  ## key* string
  ## The name of the StatKey or Property to which the condition applies
  ##
  ## operator* string
  ## Comparison operator to use.
  ## Default value is EXISTS, i.e. checks the existence of stat or property.
  ## Enum:
  ## [ EQ, NOT_EQ, LIKE, LT, GT, LT_EQ, GT_EQ, IN, NOT_IN, EXISTS, CONTAINS, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH, NOT_ENDS_WITH, NOT_CONTAINS, REGEX, NOT_REGEX,   ## NOT_EXISTS, EMPTY, NOT_EMPTY ]
  ##
  ## stringValue string
  ## String value to which we need to compare to
  #   [[inputs.vrops.tanzu_project_property_conditions.conditions]]
  #     key = "config|name"
  #     operator = "CONTAINS"
  #     stringValue = "pjp"

  ## Tanzu VM Filter
  ## Object used to lookup Tanzu VM with various filtering criteria
  ## Indicates the conjunction of the filtering criteria
  ## Either all of the filtering criteria apply together (AND operation) or any of the filtering criteria could be applied (OR operation)
  ## Defaults to OR.
  # [inputs.vrops.tanzu_vm_property_conditions]
  #   conjunctionOperator="AND"

  ## key* string
  ## The name of the StatKey or Property to which the condition applies
  ##
  ## operator* string
  ## Comparison operator to use.
  ## Default value is EXISTS, i.e. checks the existence of stat or property.
  ## Enum:
  ## [ EQ, NOT_EQ, LIKE, LT, GT, LT_EQ, GT_EQ, IN, NOT_IN, EXISTS, CONTAINS, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH, NOT_ENDS_WITH, NOT_CONTAINS, REGEX, NOT_REGEX,   ## NOT_EXISTS, EMPTY, NOT_EMPTY ]
  ##
  ## stringValue string
  ## String value to which we need to compare to
  #   [[inputs.vrops.tanzu_vm_property_conditions.conditions]]
  #     key = "config|name"
  #     operator = "NOT_CONTAINS"
  #     stringValue = "-control-plane-"

  ## Statsd data translation templates, more info can be read here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/TEMPLATE_PATTERN.md
  metric_separator = "_"
  templates = [
      "measurement.field*",
      "cloudzones.* measurement.cloudzone.cloudzone.field*",
      "net.*.* measurement.net.field*",
      "virtualDisk.*.* measurement.net.field*",
      "virtualDisk.*.*.* measurement.net.net.field*",
      "guestfilesystem.*.* measurement.device.field*",
      "guestfilesystem.*.*.* measurement.device.device.field*"
  ]

  ## Amount of time allowed to complete the HTTP(s) request.
  # timeout = "5s"

  ## Optional TLS Config
  # tls_ca = /path/to/cafile
  # tls_cert = /path/to/certfile
  # tls_key = /path/to/keyfile
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
```

## vROps -InfluxDB line-protocol Templates

The plugin supports specifying templates for transforming vROps Stat Keys into
InfluxDB measurement names and tags. The templates have a _measurement_ keyword,
which can be used to specify parts of the key that are to be used in the
measurement name. Other words in the template are used as tag names. For
example, the following template:

```toml
templates = [
    "measurement.field*"
]
```

would result in the following transformation:

```shell
summary.vcpu_allocated 30.0
=summary,vcpu_allocatedd=30.0
```

Users can also filter the template to use based on the key of the stat,
using glob matching, like so:

```toml
templates = [
    "cloudzones.* measurement.cloudzone.cloudzone.field*"
]
```

which would result in the following transformation:

```shell
cloudzones.tst-mgmt-vc.Datacenter.vCPUAllocated 4.0
=cloudzones,cloudzones=tst-mgmt-vc_Datacenter vCPUAllocated=4.0
```

Consult the [Template Patterns](/docs/TEMPLATE_PATTERN.md) documentation for
additional details.

## Metrics

- vrops_project_badge
  - compliance
  - efficiency
  - health
  - risk
- vrops_project_capacity
  - depReclaimableVCpu
  - reclaimableVCpuPerc
  - depTotalVCpu
  - depReclaimableMem
  - reclaimableMemPerc
  - depTotalMem
  - depReclaimableDiskSpace
  - reclaimableDiskSpacePerc
  - depTotalDiskSpace
- vrops_project_cloudzones
  - vCPUAllocated
  - vCPULimit
  - vCPUUtilized
  - MemoryAllocated
  - MemoryLimit
  - MemoryUtilized
  - storageAllocated
  - StorageLimit
  - StorageUtilized
- vrops_project_summary
  - vcpu_allocated
  - memory_allocated
  - storage_allocated
  - TotalBlueprints
  - TotalCloudZones
  - TotalDeployments
  - TotalPublicCloudZones
  - VMCount
  - metering_cpu
  - metering_memory
  - metering_storage
  - metering_partialPrice
  - metering_value
  - metering_additional
- vrops_project_system_attributes
  - alert_count_critical
  - alert_count_immediate
  - all_metrics
  - availability
  - child_all_metrics
  - health
  - total_alarms
  - alert_count_info
  - alert_count_warning
  - self_alert_count
  - total_alert_count
- vrops_vm_sys
  - osUptime_latest
  - poweredOn
- vrops_vm_cpu
  - usage_average
  - workload
  - iowaitPct
  - peak_vcpu_usage
  - peak_vcpu_ready
  - readyPct
  - costopPct
  - swapwaitPct
  - vcpu_usage_disparity
  - vm_capacity_provisioned
  - effective_limit
  - corecount_provisioned
  - demandPct
  - demandmhz
  - capacity_contentionPct
  - run_summation
  - overlap_summation
- vrops_vm_config
  - hardware_num_Cpu
  - hardware_disk_Space
- vrops_vm_mem-host
  - workload
- vrops_vm_mem
  - usage_average
  - workload
  - guest_demand
  - swapped_average
  - swapoutRate_average
  - swapinRate_average
  - compressed_average
  - guest_provisioned
  - overhead_average
  - reservation_used
  - effective_limit
- vrops_vm_diskspace
  - provisionedSpace
  - used
  - snapshot
  - snapshot_used
  - notshared
  - activeNotShared
  - workload
- vrops_vm_guest
  - cpu_queue
  - contextSwapRate_latest
  - disk_queue
  - used_memory
  - mem_free_latest
  - mem_physUsable_latest
  - mem.needed_latest
  - page_outRate_latest
  - page_inRate_latest
  - page_size_latest
  - swap_spaceRemaining_latest
  - tools_running_status
- vrops_vm_guestfilesystem
  - capacity_total
  - freespace_total
  - usage_total
  - percentage_total
- vrops_vm_net
  - transmitted_average
  - received_average
  - usage_average
  - droppedTx_summation
  - broadcastTx_summation
  - multicastTx_summation
  - droppedTx_summation_sum
  - multicastTx_summation_sum
  - broadcastTx_summation_sum
- vrops_vm_storage
  - totalReadLatency_average
  - totalWriteLatency_average
- vrops_vm_virtualDisk
  - read_average
  - write_average
  - totalLatency
  - totalReadLatency_average
  - totalWriteLatency_average
  - numberReadAveraged_average
  - numberWriteAveraged_average
  - commandsAveraged_average
- vrops_vm_OnlineCapacityAnalytics
  - capacityRemainingPercentage
  - timeRemaining
  - cpu_capacityRemaining
  - cpu_recommendedSize
  - cpu_timeRemaining
  - mem_capacityRemaining
  - mem_recommendedSize
  - mem_timeRemaining
  - diskspace_capacityRemaining
  - diskspace_recommendedSize
  - diskspace_timeRemaining
- vrops_tanzu_project_badge
  - compliance
  - efficiency
  - health
  - risk
  - workload
- vrops_tanzu_project_cpu
  - capacity_contentionPct
  - demandmhz
  - dynamic_entitlement
  - effective_limit
  - estimated_entitlement
  - reservation_used
  - usagemhz_average
  - workload
- vrops_tanzu_project_mem
  - active_average
  - consumed_average
  - dynamic_entitlement
  - effective_limit
  - granted_average
  - guest_demand
  - guest_provisioned
  - guest_usage
  - host_contentionPct
  - overhead_average
  - reservation_used
  - shared_average
  - swapinRate_average
  - swapoutRate_average
  - usage_average
  - workload
- vrops_tanzu_project_onlinecapacityanalytics
  - capacityRemainingPercentage
  - cpu_capacityRemaining
  - cpu_recommendedSize
  - cpu_timeRemaining
  - mem_capacityRemaining
  - mem_recommendedSize
  - mem_timeRemaining
  - timeRemaining
- vrops_tanzu_project_summary
  - number_running_vms
  - number_vm_templates
  - total_number_vms
- vrops_tanzu_project_systemattributes
  - alert_count_critical
  - alert_count_immediate
  - alert_count_info
  - alert_count_warning
  - all_metrics
  - availability
  - child_all_metrics
  - health
  - self_alert_count
  - total_alarms
  - total_alert_count
- vrops_tanzu_vm_badge
  - compliance
  - efficiency
  - health
  - risk
  - workload
- vrops_tanzu_vm_config
  - hardware_disk_Space
  - hardware_num_Cpu
- vrops_tanzu_vm_cpu
  - capacity_contentionPct
  - corecount_provisioned
  - costopPct
  - demandmhz
  - demandPct
  - effective_limit
  - iowaitPct
  - overlap_summation
  - peak_vcpu_ready
  - peak_vcpu_usage
  - readyPct
  - run_summation
  - swapwaitPct
  - usage_average
  - usagemhz_average
  - usagemhz_average_daily
  - vcpu_usage_disparity
  - vm_capacity_provisioned
  - workload
- vrops_tanzu_vm_diskspace
  - perDsUsed
  - activeNotShared
  - notshared
  - provisionedSpace
  - snapshot
  - snapshot_used
  - used
  - workload
- vrops_tanzu_vm_diskspace-total
  - workload
- vrops_tanzu_vm_guest
  - contextSwapRate_latest
  - cpu_queue
  - disk_queue
  - mem_free_latest
  - mem_needed_latest
  - mem_physUsable_latest
  - page_inRate_latest
  - page_outRate_latest
  - page_size_latest
  - swap_spaceRemaining_latest
  - tools_running_status
  - used_memory
- vrops_tanzu_vm_guestfilesystem
  - capacity
  - percentage
  - usage
  - capacity_total
  - percentage_total
  - usage_total
- vrops_tanzu_vm_mem
  - balloonPct
  - compressed_average
  - consumed_average
  - consumed_average_daily
  - consumedPct
  - effective_limit
  - guest_demand
  - guest_provisioned
  - guest_usage
  - host_contentionPct
  - host_demand
  - nonzero_active
  - overhead_average
  - overheadMax_average
  - reservation_used
  - swapinRate_average
  - swapoutRate_average
  - swapped_average
  - usage_average
  - vmMemoryDemand
  - workload
- vrops_tanzu_vm_mem-host
  - workload
- vrops_tanzu_vm_net
  - droppedPct
  - packetsRxPerSec
  - packetsTxPerSec
  - broadcastTx_summation
  - broadcastTx_summation_sum
  - droppedTx_summation
  - droppedTx_summation_sum
  - multicastTx_summation
  - multicastTx_summation_sum
  - received_average
  - transmitted_average
  - usage_average
- vrops_tanzu_vm_onlinecapacityanalytics
  - capacityRemainingPercentage
  - cpu_capacityRemaining
  - cpu_recommendedSize
  - cpu_timeRemaining
  - diskspace_capacityRemaining
  - diskspace_recommendedSize
  - diskspace_timeRemaining
  - mem_capacityRemaining
  - mem_recommendedSize
  - mem_timeRemaining
  - timeRemaining
- vrops_tanzu_vm_performance
  - number_of_kpis_breached
- vrops_tanzu_vm_power
  - energy_summation_sum
- vrops_tanzu_vm_rescpu
  - actav1_latest
  - actav5_latest
- vrops_tanzu_vm_storage
  - totalReadLatency_average
  - totalWriteLatency_average
- vrops_tanzu_vm_summary
  - idle
  - oversized
  - oversized_memory
  - oversized_vcpus
  - poweredOff
  - running
  - snapshotSpace
  - undersized
  - undersized_memory
  - undersized_vcpus
- vrops_tanzu_vm_sys
  - osUptime_latest
  - poweredOn
- vrops_tanzu_vm_systemattributes
  - alert_count_critical
  - alert_count_immediate
  - alert_count_info
  - alert_count_warning
  - all_metrics
  - availability
  - child_all_metrics
  - health
  - self_alert_count
  - total_alarms
  - total_alert_count
- vrops_tanzu_vm_virtualdisk
  - commandsAveraged_average
  - numberReadAveraged_average
  - numberWriteAveraged_average
  - totalLatency
  - totalReadLatency_average
  - totalWriteLatency_average
  - usage
  - vDiskOIO
  - numberReadAveraged_average
  - numberWriteAveraged_average
  - read_average
  - totalReadLatency_average
  - totalWriteLatency_average
  - write_average
  - peak_vDisk_iops
  - peak_vDisk_readLatency
  - peak_vDisk_writeLatency

For more information about the metrics, refer to the [documentation][vrops-info].

[vrops-info]: https://seversky.atlassian.net/wiki/spaces/CSHD/pages/2069102615/vROps

## Example Output

```shell
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC vCPUUtilized=33.33333333333333 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC alert_count_warning=2 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC all_metrics=32 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC depReclaimableMem=0 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC depTotalDiskSpace=391.1988525390625 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC vCPULimit=24 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC metering_additional=493.1822204589844 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC depTotalVCpu=8 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC MemoryUtilized=4.096 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC TotalCloudZones=1 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC VMCount=4 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC TotalDeployments=4 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC depReclaimableDiskSpace=0 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC reclaimableDiskSpacePerc=0 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC vCPUAllocated=8 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC reclaimableMemPerc=0 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC StorageLimit=1000000 1678077167000000000
vrops_project_badge,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC risk=0 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC storageAllocated=400 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC MemoryAllocated=8388608 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC StorageUtilized=0.04 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC depTotalMem=8388608 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC total_alert_count=2 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC vcpu_allocated=33.33333333333333 1678077167000000000
vrops_project_badge,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC compliance=-1 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC alert_count_info=0 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC alert_count_critical=0 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC TotalBlueprints=0 1678077167000000000
vrops_project_cloudzones,cloudzone=tst-mgmt-vc_Datacenter,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC MemoryLimit=204800000 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC metering_partialPrice="false" 1678077167000000000
vrops_project_badge,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC efficiency=100 1678077167000000000
vrops_project_badge,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC health=100 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC total_alarms=112 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC metering_value=839.8155717849731 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC health=76 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC metering_memory=203.16827392578125 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC metering_storage=98.67170333862305 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC reclaimableVCpuPerc=0 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC metering_cpu=44.79337406158447 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC TotalPublicCloudZones=0 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC alert_count_immediate=0 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC self_alert_count=0 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC availability=1 1678077167000000000
vrops_project_systemattributes,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC child_all_metrics=6752 1678077167000000000
vrops_project_capacity,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC depReclaimableVCpu=0 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC memory_allocated=4.096 1678077167000000000
vrops_project_summary,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC storage_allocated=0.04 1678077167000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic5,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_cpu,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0.392 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic1,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0 1678077893000000000
vrops_vm_virtualdisk,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,device=scsi0_0,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 read_average=0 1678077893000000000
vrops_vm_diskspace,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 provisionedSpace=202.30859375 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic1,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_virtualdisk,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,device=scsi0_0,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 write_average=8.6 1678077893000000000
vrops_vm_cpu,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 corecount_provisioned=2 1678077893000000000
vrops_vm_guestfilesystem,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 percentage_total=12.917854895043252 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic2,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=4000,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic2,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=4000,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic0,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic4,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0 1678077893000000000
vrops_vm_sys,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 poweredOn=1 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic5,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_guestfilesystem,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 freespace_total=85.16608810424805 1678077893000000000
vrops_vm_virtualdisk,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 read_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic3,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic2,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0 1678077893000000000
vrops_vm_config,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 hardware_disk_Space=100 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic3,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_mem,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=26.756149927775063 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vusb0,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_virtualdisk,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 write_average=8.6 1678077893000000000
vrops_vm_guestfilesystem,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_total=12.633625030517578 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=4000,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic5,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 usage_average=0 1678077893000000000
vrops_vm_mem,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 guest_provisioned=2097152 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic1,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vmnic4,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 received_average=0 1678077893000000000
vrops_vm_net,deployment_id=27e7d17f-2343-4c9e-9d13-e4e4e1ebb616,deployment_name=kor-uni-03,host=S2100113,net=vusb0,tenant_id=1a099fbd-379e-4153-83bf-19739de99b69,tenant_name=KOR-RC,vm_id=f2ecbace-9f7b-4a73-a97a-e512c3fe4040,vm_name=korlab.121 transmitted_average=0 1678077893000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 alert_count_warning=0 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 all_metrics=48 1682408481000000000
vrops_tanzu_project_summary,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 total_number_vms=3 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 cpu_capacityRemaining=488365.2701706162 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 guest_usage=20951292.8 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 swapinRate_average=0 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 usagemhz_average=2953.8 1682408481000000000
vrops_tanzu_project_badge,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 workload=38.253067070874735 1682408481000000000
vrops_tanzu_project_supermetric,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 sm_f61a7d58-8309-4921-810b-e13a8a002b5a=254 1682408481000000000
vrops_tanzu_project_badge,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 risk=0 1682408481000000000
vrops_tanzu_project_summary,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 number_running_vms=3 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 reservation_used=0 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 total_alert_count=0 1682408481000000000
vrops_tanzu_project_badge,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 compliance=-1 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 alert_count_info=0 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 capacityRemainingPercentage=47.94168158167578 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 consumed_average=20951292.8 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 timeRemaining=366 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 workload=0.6148372714029757 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 active_average=7253343.466666667 1682408481000000000
vrops_tanzu_project_badge,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 health=100 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 health=100 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 usage_average=34.58663686116537 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 guest_demand=7978677.813333333 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 alert_count_immediate=0 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 guest_provisioned=20971520 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 overhead_average=172807.46666666667 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 reservation_used=283648 1682408481000000000
vrops_tanzu_project_supermetric,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 sm_d2ef43c1-a365-406a-b871-54e9566916e7=8 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 child_all_metrics=3768 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 effective_limit=497088 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 mem_recommendedSize=10858124.118365692 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 granted_average=20971520 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 mem_timeRemaining=366 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 demandmhz=3034.6666666666665 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 mem_capacityRemaining=9999491.817503063 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 workload=38.253067070874735 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 dynamic_entitlement=20857615.935868755 1682408481000000000
vrops_tanzu_project_supermetric,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 sm_c98f08a0-c4af-4301-a2f4-6e97b327831d=16 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 swapoutRate_average=0 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 alert_count_critical=0 1682408481000000000
vrops_tanzu_project_summary,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 number_vm_templates=0 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 shared_average=20867.2 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 cpu_timeRemaining=366 1682408481000000000
vrops_tanzu_project_badge,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 efficiency=100 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 capacity_contentionPct=0.03773333333333333 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 total_alarms=16 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 self_alert_count=0 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 estimated_entitlement=493572.3333333333 1682408481000000000
vrops_tanzu_project_systemattributes,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 availability=1 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 effective_limit=3910576128 1682408481000000000
vrops_tanzu_project_supermetric,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 sm_fa9d590d-35ce-4706-8384-d62927554b88=2 1682408481000000000
vrops_tanzu_project_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 host_contentionPct=0 1682408481000000000
vrops_tanzu_project_onlinecapacityanalytics,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 cpu_recommendedSize=246786.16666666663 1682408481000000000
vrops_tanzu_project_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001 dynamic_entitlement=27930 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=Aggregateofallinstances,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn totalReadLatency_average=0 1682408481000000000
vrops_tanzu_vm_virtualdisk,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn peak_vDisk_iops=60.266666666666666 1682408481000000000
vrops_tanzu_vm_guestfilesystem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn capacity_total=689.2849884033203 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_2,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn write_average=0 1682408481000000000
vrops_tanzu_vm_diskspace,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn used=126.27726807352155 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_2,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn read_average=0 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_1,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn read_average=0 1682408481000000000
vrops_tanzu_vm_config,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn hardware_num_Cpu=4 1682408481000000000
vrops_tanzu_vm_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn usage_average=8.253333333333334 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_0,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn read_average=435.06666666666666 1682408481000000000
vrops_tanzu_vm_diskspace,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn provisionedSpace=412.35546875 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_0,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn write_average=904.2666666666667 1682408481000000000
vrops_tanzu_vm_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn corecount_provisioned=4 1682408481000000000
vrops_tanzu_vm_guestfilesystem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn percentage_total=43.00272102919221 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=Aggregateofallinstances,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn totalLatency=0.2117994100294985 1682408481000000000
vrops_tanzu_vm_sys,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn poweredOn=1 1682408481000000000
vrops_tanzu_vm_virtualdisk,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn read_average=435.06666666666666 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_1,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn write_average=15.466666666666667 1682408481000000000
vrops_tanzu_vm_config,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn hardware_disk_Space=202 1682408481000000000
vrops_tanzu_vm_net,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn transmitted_average=27.533333333333335 1682408481000000000
vrops_tanzu_vm_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn usage_average=42.587264378865555 1682408481000000000
vrops_tanzu_vm_net,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn usage_average=118 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=Aggregateofallinstances,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn totalWriteLatency_average=0.26666666666666666 1682408481000000000
vrops_tanzu_vm_virtualdisk,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn write_average=920 1682408481000000000
vrops_tanzu_vm_guestfilesystem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn usage_total=296.4113006591797 1682408481000000000
vrops_tanzu_vm_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn guest_provisioned=8388608 1682408481000000000
vrops_tanzu_vm_net,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=ad6f5c67-6557-4cc1-9d23-fe386fccbf90,vm_name=pjp-dev001-md-0-5bb975b56c-v7bbn received_average=90 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=Aggregateofallinstances,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 totalReadLatency_average=0 1682408481000000000
vrops_tanzu_vm_virtualdisk,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 peak_vDisk_iops=62.93333333333334 1682408481000000000
vrops_tanzu_vm_guestfilesystem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 capacity_total=298.03282928466797 1682408481000000000
vrops_tanzu_vm_diskspace,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 used=96.62474834360182 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_1,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 read_average=0 1682408481000000000
vrops_tanzu_vm_config,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 hardware_num_Cpu=4 1682408481000000000
vrops_tanzu_vm_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 usage_average=7.159333333333333 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_0,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 read_average=9.666666666666666 1682408481000000000
vrops_tanzu_vm_diskspace,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 provisionedSpace=112.33203125 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_0,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 write_average=1011.9333333333333 1682408481000000000
vrops_tanzu_vm_cpu,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 corecount_provisioned=4 1682408481000000000
vrops_tanzu_vm_guestfilesystem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 percentage_total=68.52583490493068 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=Aggregateofallinstances,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 totalLatency=0 1682408481000000000
vrops_tanzu_vm_sys,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 poweredOn=1 1682408481000000000
vrops_tanzu_vm_virtualdisk,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 read_average=9.666666666666666 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=scsi0_1,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 write_average=0 1682408481000000000
vrops_tanzu_vm_config,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 hardware_disk_Space=52 1682408481000000000
vrops_tanzu_vm_net,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 transmitted_average=38 1682408481000000000
vrops_tanzu_vm_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 usage_average=30.13167381286621 1682408481000000000
vrops_tanzu_vm_net,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 usage_average=70.93333333333334 1682408481000000000
vrops_tanzu_vm_virtualdisk,device=Aggregateofallinstances,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 totalWriteLatency_average=0 1682408481000000000
vrops_tanzu_vm_virtualdisk,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 write_average=1011.9333333333333 1682408481000000000
vrops_tanzu_vm_guestfilesystem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 usage_total=204.22948455810547 1682408481000000000
vrops_tanzu_vm_mem,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 guest_provisioned=8388608 1682408481000000000
vrops_tanzu_vm_net,host=S2100113,tenant_id=e10d6ccd-6b47-4811-be59-924f9872cf15,tenant_name=PJP-DEV001,vm_id=061fc8b8-e2d6-40cc-834b-d0e3fd4994a0,vm_name=pjp-dev001-md-0-5bb975b56c-fqlh7 received_average=32.46666666666667 1682408481000000000
```
