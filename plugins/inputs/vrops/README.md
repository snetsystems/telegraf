# VMware vROps(vRealize Operations) Input Plugin

The VMware vROps plugin uses the vROps API to gather metrics.

- Projects
- VMs

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
  ## "projects", "vms"
  enabled_services = ["projects", "vms"]

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
```
