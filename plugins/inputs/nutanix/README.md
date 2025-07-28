# Nutanix PRISM Input Plugin

The Nutanix PRISM plugin uses the Nutanix API to gather metrics.

- Cluster
- Hosts
- Disks
- VMs

## Configuration

```toml @sample.conf
# Collects performance metrics from Nutanix PRISM services
[[inputs.nutanix]]
  ## List of Nutanix URL to be monitored. These three lines must be uncommented
  ## and edited for the plugin to work.
  url = "https://localhost:9440/api/nutanix/v3/groups"
  username = "nutanix"
  password = "nutanix"

  ## Available services are:
  ## "cluster", "hosts", "disks", "vms"
  enabled_services = ["cluster", "hosts", "disks", "vms"]

  ## Amount of time allowed to complete the HTTP(s) request.
  # timeout = "5s"

  ## Optional TLS Config
  # tls_ca = /path/to/cafile
  # tls_cert = /path/to/certfile
  # tls_key = /path/to/keyfile
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
```

## Metrics

- nutanix_cluster
  - num_nodes
  - num_vms
  - version
  - external_ip_address
  - hypervisor_cpu_usage_ppm
  - overall_memory_usage_ppm
  - aggregate_hypervisor_memory_usage_ppm
  - hypervisor_cpu_usage_percent
  - overall_memory_usage_percent
  - aggregate_hypervisor_memory_usage_percent
  - controller_num_iops
  - controller_num_read_iops
  - controller_num_write_iops
  - controller_io_bandwidth_kBps
  - controller_read_io_bandwidth_kBps
  - controller_write_io_bandwidth_kBps
  - controller_avg_io_latency_usecs
  - controller_avg_read_io_latency_usecs
  - controller_avg_write_io_latency_usecs
  - power_consumption_instant_watt
- nutanix_hosts
  - num_vms
  - hypervisor_full_name
  - service_vm_ipv4_address
  - block_model_name
  - hypervisor_cpu_usage_ppm
  - overall_memory_usage_ppm
  - aggregate_hypervisor_memory_usage_ppm
  - hypervisor_cpu_usage_percent
  - overall_memory_usage_percent
  - aggregate_hypervisor_memory_usage_percent
  - controller_num_iops
  - controller_num_read_iops
  - controller_num_write_iops
  - controller_io_bandwidth_kBps
  - controller_read_io_bandwidth_kBps
  - controller_write_io_bandwidth_kBps
  - controller_avg_io_latency_usecs
  - controller_avg_read_io_latency_usecs
  - controller_avg_write_io_latency_usecs
  - power_consumption_instant_watt
- nutanix_disks
  - num_iops
  - num_read_iops
  - num_write_iops
  - avg_io_latency_usecs
  - avg_read_io_latency_usecs
  - avg_write_io_latency_usecs
  - io_bandwidth_kBps
  - read_io_bandwidth_kBps
  - write_io_bandwidth_kBps
  - storage_usage_ppm
  - storage_usage_percent
- nutanix_vms
  - gpus_in_use
  - power_state
  - hypervisor_cpu_usage_ppm
  - hypervisor.cpu_ready_time_ppm
  - memory_usage_ppm
  - disk_usage_ppm
  - hypervisor_cpu_usage_percent
  - hypervisor.cpu_ready_time_percent
  - memory_usage_percent
  - disk_usage_percent
  - controller.snapshot_usage_bytes
  - controller_user_bytes
  - controller_io_bandwidth_kBps
  - controller_read_io_bandwidth_kBps
  - controller_write_io_bandwidth_kBps
  - controller_avg_io_latency_usecs
  - controller_avg_read_io_latency_usecs
  - controller_avg_write_io_latency_usecs
  - controller_num_iops
  - controller_num_read_iops
  - controller_num_write_iops
  - hypervisor_num_received_bytes
  - hypervisor_num_receive_packets_dropped
  - hypervisor_num_transmitted_bytes
  - hypervisor_num_transmit_packets_dropped
  - controller.wss_3600s_read_MB
  - controller.wss_3600s_write_MB
  - controller.wss_3600s_union_MB
  - controller.shared_usage_bytes

For more information about the metrics, refer to the Cloudhub Document.

## Example Output

```shell
nutanix_cluster,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems aggregate_hypervisor_memory_usage_ppm=82018i,controller_avg_io_latency_usecs=1030i,controller_avg_read_io_latency_usecs=403i,controller_avg_write_io_latency_usecs=1031i,controller_io_bandwidth_kBps=580i,controller_num_iops=61i,controller_num_read_iops=0i,controller_num_write_iops=61i,controller_read_io_bandwidth_kBps=0i,controller_write_io_bandwidth_kBps=580i,hypervisor_cpu_usage_ppm=50541i,num_nodes=3i,num_vms=8i,overall_memory_usage_ppm=96965i 1750050343000000000
nutanix_hosts,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1 aggregate_hypervisor_memory_usage_ppm=76543i,controller_avg_io_latency_usecs=0i,controller_avg_read_io_latency_usecs=0i,controller_avg_write_io_latency_usecs=0i,controller_io_bandwidth_kBps=0i,controller_num_iops=0i,controller_num_read_iops=0i,controller_num_write_iops=0i,controller_read_io_bandwidth_kBps=0i,controller_write_io_bandwidth_kBps=0i,hypervisor_cpu_usage_ppm=51143i,overall_memory_usage_ppm=109099i 1750050343000000000
nutanix_hosts,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2 aggregate_hypervisor_memory_usage_ppm=96599i,controller_avg_io_latency_usecs=835i,controller_avg_read_io_latency_usecs=436i,controller_avg_write_io_latency_usecs=835i,controller_io_bandwidth_kBps=820i,controller_num_iops=77i,controller_num_read_iops=0i,controller_num_write_iops=77i,controller_read_io_bandwidth_kBps=0i,controller_write_io_bandwidth_kBps=819i,hypervisor_cpu_usage_ppm=84350i,overall_memory_usage_ppm=96599i 1750050343000000000
nutanix_hosts,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3 aggregate_hypervisor_memory_usage_ppm=72912i,controller_avg_io_latency_usecs=1596i,controller_avg_read_io_latency_usecs=0i,controller_avg_write_io_latency_usecs=1596i,controller_io_bandwidth_kBps=0i,controller_num_iops=0i,controller_num_read_iops=0i,controller_num_write_iops=0i,controller_read_io_bandwidth_kBps=0i,controller_write_io_bandwidth_kBps=0i,hypervisor_cpu_usage_ppm=38328i,overall_memory_usage_ppm=85196i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=04ea104d-a8c0-46b0-9201-dcf2bde1875b,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHAB226402RH3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=05878624-ab44-48bb-a9df-6847a4d3b773,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHKE2155008F375FGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="DAS-SATA",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=0b350ff6-58f6-414a-b0c0-e9e2fac2a728,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHLJ212400K81P0FGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="DAS-SATA",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=19c49fcd-8248-4c75-a6a4-613065d76be1,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHAB226402KK3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=1fa7f395-8abc-4fa1-b9f4-b6dc60ed8a2e,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHAB2506000Q3P8EGN avg_io_latency_usecs=130i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=296f2c45-4b43-4b53-be93-9247afd342ae,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHAB226402LV3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=3be529cb-2153-4e6c-9d9f-80ae4def1a0b,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHLJ212400NQ1P0FGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="DAS-SATA",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=4103001c-24e3-4bcd-a29e-7fa19f4af559,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHAB250600B23P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=5710df54-3987-45c1-9ab4-23879729b9a4,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHAB250600GT3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=6716e1e7-fdcd-4961-9a51-73ce4ec9b264,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHKE2155008X375FGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="DAS-SATA",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=73fef6c8-147b-4ce1-bf9d-35e8b82e0658,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHAB2506009R3P8EGN avg_io_latency_usecs=261i,io_bandwidth_kBps=34i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=34i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=74f7f587-ddc8-4b38-948d-6553639e41dd,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHAB2506002R3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=75bb2194-8831-4986-856e-e9963e77c1e2,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHLJ212400P51P0FGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="DAS-SATA",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=7970aa72-acd6-44bc-af18-3fd405a562bd,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHAB2506004N3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=79e7cb15-6193-42c2-a817-9723c5740886,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHAB2506003V3P8EGN avg_io_latency_usecs=226i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=7addc18c-58ad-45ed-b5ac-286345c75103,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHAB226402J93P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=89db3f3b-34bd-4bbf-b72a-d3f5ea2b7dea,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHAB250600G53P8EGN avg_io_latency_usecs=69i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=a7e9aec6-c7f9-4400-8040-619f2edc0f48,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHAB2506004J3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=b73ec174-d773-44f3-b2bc-aee04b9b93ad,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,serial=PHAB250600A73P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=c3a29992-feab-4f8e-afdb-7dfd0769c885,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHAB250300973P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=c988ef42-2652-4f0a-aeef-d5778f53cb47,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHAB226402KX3P8EGN avg_io_latency_usecs=16i,io_bandwidth_kBps=1i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=1i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=cd11b4d8-15ee-4d63-8f2f-e2ddfd388e92,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHAB226402KU3P8EGN avg_io_latency_usecs=20i,io_bandwidth_kBps=1i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=1i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=ddd80f7f-1d38-45b7-ba6d-97f669305926,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,serial=PHAB2506002T3P8EGN avg_io_latency_usecs=0i,io_bandwidth_kBps=0i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="SSD-PCIe",write_io_bandwidth_kBps=0i 1750050343000000000
nutanix_disks,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,disk_id=f1ccda17-ed1d-4178-aeb5-ce0cc14b7c7a,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,serial=PHKE21550013375FGN avg_io_latency_usecs=140i,io_bandwidth_kBps=2i,num_iops=0i,num_read_iops=0i,num_write_iops=0i,read_io_bandwidth_kBps=0i,storage_tier="DAS-SATA",write_io_bandwidth_kBps=2i 1750050343000000000
nutanix_vms,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=34679686-a81b-46a7-848d-222222222224,node_name=ahv2,project_name=_internal,vm_id=922158db-c739-4c01-9732-96a630f5165b,vm_name=PC-1 controller_avg_io_latency_usecs=835i,controller_avg_read_io_latency_usecs=436i,controller_avg_write_io_latency_usecs=835i,controller_io_bandwidth_kBps=820i,controller_num_iops=77i,controller_num_read_iops=0i,controller_num_write_iops=77i,controller_read_io_bandwidth_kBps=0i,controller_shared_usage_bytes=0i,controller_snapshot_usage_bytes=0i,controller_user_bytes=81538619904i,controller_write_io_bandwidth_kBps=819i,controller_wss_3600s_read_MB=734i,controller_wss_3600s_union_MB=14296i,controller_wss_3600s_write_MB=14074i,disk_usage_ppm=71542i,gpus_in_use=0i,hypervisor_cpu_ready_time_ppm=15i,hypervisor_cpu_usage_ppm=340269i,hypervisor_num_receive_packets_dropped=0i,hypervisor_num_received_bytes=149323i,hypervisor_num_transmitted_bytes=1410663i,memory_usage_ppm=661419i,power_state="on",project_reference="30fe5fb6-d394-4055-92f8-625507b3797c" 1750050343000000000
nutanix_vms,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,project_name=_internal,vm_id=7bc3b069-97f7-4f35-b398-74652598d3ca,vm_name=jm-test controller_avg_io_latency_usecs=0i,controller_avg_read_io_latency_usecs=0i,controller_avg_write_io_latency_usecs=0i,controller_io_bandwidth_kBps=0i,controller_num_iops=0i,controller_num_read_iops=0i,controller_num_write_iops=0i,controller_read_io_bandwidth_kBps=0i,controller_shared_usage_bytes=0i,controller_snapshot_usage_bytes=0i,controller_user_bytes=0i,controller_write_io_bandwidth_kBps=0i,controller_wss_3600s_read_MB=1i,controller_wss_3600s_union_MB=1i,controller_wss_3600s_write_MB=0i,disk_usage_ppm=0i,gpus_in_use=0i,hypervisor_cpu_ready_time_ppm=6i,hypervisor_cpu_usage_ppm=15638i,hypervisor_num_receive_packets_dropped=0i,hypervisor_num_received_bytes=389i,hypervisor_num_transmitted_bytes=0i,power_state="on",project_reference="30fe5fb6-d394-4055-92f8-625507b3797c" 1750050343000000000
nutanix_vms,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=03e5b997-f3c3-4304-ab8c-111111111108,node_name=ahv1,project_name=_internal,vm_id=c3199ca5-4bf7-44d4-a1d0-943494a53226,vm_name=jm-test-1 controller_avg_io_latency_usecs=0i,controller_avg_read_io_latency_usecs=0i,controller_avg_write_io_latency_usecs=0i,controller_io_bandwidth_kBps=0i,controller_num_iops=0i,controller_num_read_iops=0i,controller_num_write_iops=0i,controller_read_io_bandwidth_kBps=0i,controller_shared_usage_bytes=0i,controller_snapshot_usage_bytes=0i,controller_user_bytes=0i,controller_write_io_bandwidth_kBps=0i,controller_wss_3600s_read_MB=1i,controller_wss_3600s_union_MB=1i,controller_wss_3600s_write_MB=0i,disk_usage_ppm=0i,gpus_in_use=0i,hypervisor_cpu_ready_time_ppm=16i,hypervisor_cpu_usage_ppm=522407i,hypervisor_num_receive_packets_dropped=0i,hypervisor_num_received_bytes=0i,hypervisor_num_transmitted_bytes=389i,power_state="on",project_reference="30fe5fb6-d394-4055-92f8-625507b3797c" 1750050343000000000
nutanix_vms,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,project_name=_internal,vm_id=42ac8c01-88c5-45a9-8d32-de6b4ec69035,vm_name=jm-dns controller_shared_usage_bytes=2018066432i,controller_snapshot_usage_bytes=0i,controller_user_bytes=2109505536i,controller_wss_3600s_read_MB=0i,controller_wss_3600s_union_MB=0i,controller_wss_3600s_write_MB=0i,disk_usage_ppm=49115i,gpus_in_use=0i,hypervisor_cpu_usage_ppm=0i,hypervisor_num_receive_packets_dropped=0i,hypervisor_num_received_bytes=0i,hypervisor_num_transmitted_bytes=0i,memory_usage_ppm=0i,power_state="off",project_reference="30fe5fb6-d394-4055-92f8-625507b3797c" 1750050343000000000
nutanix_vms,cluster_id=000630c3-4970-c20c-3800-0025b588888f,cluster_name=cluster,host=snetsystems,node_id=e2b4c99d-0d98-40a7-9595-3333333333b6,node_name=ahv3,project_name=_internal,vm_id=0ca0fffc-ff24-47a2-ae3a-1f7f4c6d73da,vm_name=jm-test2 controller_avg_io_latency_usecs=1596i,controller_avg_read_io_latency_usecs=0i,controller_avg_write_io_latency_usecs=1596i,controller_io_bandwidth_kBps=0i,controller_num_iops=0i,controller_num_read_iops=0i,controller_num_write_iops=0i,controller_read_io_bandwidth_kBps=0i,controller_shared_usage_bytes=2018066432i,controller_snapshot_usage_bytes=0i,controller_user_bytes=2121601024i,controller_write_io_bandwidth_kBps=0i,controller_wss_3600s_read_MB=0i,controller_wss_3600s_union_MB=13i,controller_wss_3600s_write_MB=13i,disk_usage_ppm=49397i,gpus_in_use=0i,hypervisor_cpu_ready_time_ppm=0i,hypervisor_cpu_usage_ppm=166i,hypervisor_num_receive_packets_dropped=0i,hypervisor_num_received_bytes=0i,hypervisor_num_transmitted_bytes=0i,memory_usage_ppm=87384i,power_state="on",project_reference="30fe5fb6-d394-4055-92f8-625507b3797c" 1750050343000000000
```
