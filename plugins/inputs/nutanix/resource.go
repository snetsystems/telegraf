// Package nutanix implements an nutanix input plugin for Telegraf
//
//go:generate ../../../tools/readme_config_includer/generator
package nutanix

type GroupMemberAttribute struct {
	Attribute string `json:"attribute"`
	Operation string `json:"operation,omitempty"`
}

type EntityPayload struct {
	EntityType            string                 `json:"entity_type"`
	GroupMemberAttributes []GroupMemberAttribute `json:"group_member_attributes"`
	FilterCriteria        string                 `json:"filter_criteria,omitempty"`
}

type ValuePoint struct {
	Time   int64    `json:"time"`
	Values []string `json:"values"`
}

type MetricData struct {
	Name     string       `json:"name"`
	DataType string       `json:"data_type"`
	Values   []ValuePoint `json:"values"`
}

type EntityResult struct {
	EntityID string       `json:"entity_id"`
	Data     []MetricData `json:"data"`
}

type GroupResult struct {
	EntityResults []EntityResult `json:"entity_results"`
}

type GroupResponse struct {
	GroupResults []GroupResult `json:"group_results"`
}

var metricPayloads = map[string]EntityPayload{
	"cluster": {
		EntityType: "cluster",
		GroupMemberAttributes: []GroupMemberAttribute{
			{Attribute: "cluster_name"},
			{Attribute: "num_nodes"},
			{Attribute: "num_vms"},
			{Attribute: "version"},
			{Attribute: "external_ip_address"},
			{Attribute: "hypervisor_cpu_usage_ppm", Operation: "avg"},
			{Attribute: "overall_memory_usage_ppm", Operation: "avg"},
			{Attribute: "aggregate_hypervisor_memory_usage_ppm", Operation: "avg"},
			{Attribute: "controller_num_iops", Operation: "avg"},
			{Attribute: "controller_num_read_iops", Operation: "last"},
			{Attribute: "controller_num_write_iops", Operation: "last"},
			{Attribute: "controller_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_read_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_write_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_avg_io_latency_usecs", Operation: "last"},
			{Attribute: "controller_avg_read_io_latency_usecs", Operation: "avg"},
			{Attribute: "controller_avg_write_io_latency_usecs", Operation: "avg"},
			{Attribute: "power_consumption_instant_watt", Operation: "avg"},
		},
	},
	"hosts": {
		EntityType: "host",
		GroupMemberAttributes: []GroupMemberAttribute{
			{Attribute: "cluster"},
			{Attribute: "cluster_name"},
			{Attribute: "node_name"},
			{Attribute: "num_vms"},
			{Attribute: "hypervisor_full_name"},
			{Attribute: "service_vm_ipv4_address"},
			{Attribute: "block_model_name"},
			{Attribute: "hypervisor_cpu_usage_ppm", Operation: "avg"},
			{Attribute: "overall_memory_usage_ppm", Operation: "avg"},
			{Attribute: "aggregate_hypervisor_memory_usage_ppm", Operation: "avg"},
			{Attribute: "controller_num_iops", Operation: "avg"},
			{Attribute: "controller_num_read_iops", Operation: "last"},
			{Attribute: "controller_num_write_iops", Operation: "last"},
			{Attribute: "controller_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_read_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_write_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_avg_io_latency_usecs", Operation: "last"},
			{Attribute: "controller_avg_read_io_latency_usecs", Operation: "avg"},
			{Attribute: "controller_avg_write_io_latency_usecs", Operation: "avg"},
			{Attribute: "power_consumption_instant_watt", Operation: "avg"},
		},
	},
	"disks": {
		EntityType: "disk",
		GroupMemberAttributes: []GroupMemberAttribute{
			{Attribute: "cluster"},
			{Attribute: "cluster_name"},
			{Attribute: "node"},
			{Attribute: "node_name"},
			{Attribute: "serial"},
			{Attribute: "storage_tier"},
			{Attribute: "storage.usage_ppm", Operation: "last"},
			{Attribute: "num_iops", Operation: "last"},
			{Attribute: "num_read_iops", Operation: "last"},
			{Attribute: "num_write_iops", Operation: "last"},
			{Attribute: "avg_io_latency_usecs", Operation: "last"},
			{Attribute: "avg_read_io_latency_usecs", Operation: "last"},
			{Attribute: "avg_write_io_latency_usecs", Operation: "last"},
			{Attribute: "io_bandwidth_kBps", Operation: "last"},
			{Attribute: "read_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "write_io_bandwidth_kBps", Operation: "last"},
		},
	},
	"vms": {
		EntityType: "mh_vm",
		GroupMemberAttributes: []GroupMemberAttribute{
			{Attribute: "vm_name"},
			{Attribute: "cluster"},
			{Attribute: "cluster_name"},
			{Attribute: "node"},
			{Attribute: "node_name"},
			{Attribute: "project_name"},
			{Attribute: "project_reference"},
			{Attribute: "gpus_in_use"},
			{Attribute: "power_state"},
			{Attribute: "vpc_name"},
			{Attribute: "num_vcpus"},
			{Attribute: "memory_size_bytes"},
			{Attribute: "hypervisor_cpu_usage_ppm", Operation: "avg"},
			{Attribute: "hypervisor.cpu_ready_time_ppm", Operation: "avg"},
			{Attribute: "memory_usage_ppm", Operation: "avg"},
			{Attribute: "disk_usage_ppm", Operation: "avg"},
			{Attribute: "controller.snapshot_usage_bytes", Operation: "avg"},
			{Attribute: "controller_user_bytes", Operation: "avg"},
			{Attribute: "controller_io_bandwidth_kBps", Operation: "avg"},
			{Attribute: "controller_read_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_write_io_bandwidth_kBps", Operation: "last"},
			{Attribute: "controller_avg_io_latency_usecs", Operation: "last"},
			{Attribute: "controller_avg_read_io_latency_usecs", Operation: "avg"},
			{Attribute: "controller_avg_write_io_latency_usecs", Operation: "avg"},
			{Attribute: "controller_num_iops", Operation: "avg"},
			{Attribute: "controller_num_read_iops", Operation: "avg"},
			{Attribute: "controller_num_write_iops", Operation: "avg"},
			{Attribute: "hypervisor_num_received_bytes", Operation: "sum"},
			{Attribute: "hypervisor_num_receive_packets_dropped", Operation: "sum"},
			{Attribute: "hypervisor_num_transmitted_bytes", Operation: "sum"},
			{Attribute: "hypervisor_num_transmit_packets_dropped", Operation: "sum"},
			{Attribute: "controller.wss_3600s_read_MB", Operation: "avg"},
			{Attribute: "controller.wss_3600s_write_MB", Operation: "avg"},
			{Attribute: "controller.wss_3600s_union_MB", Operation: "avg"},
			{Attribute: "controller.shared_usage_bytes", Operation: "avg"},
		},
		FilterCriteria: "is_cvm==0",
	},
}
