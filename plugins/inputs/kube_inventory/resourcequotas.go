package kube_inventory

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/influxdata/telegraf"
)

func collectResourceQuotas(ctx context.Context, acc telegraf.Accumulator, ki *KubernetesInventory) {
	list, err := ki.client.getResourceQuotas(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, i := range list.Items {
		ki.gatherResourceQuota(i, acc)
	}
}

func (ki *KubernetesInventory) gatherResourceQuota(r corev1.ResourceQuota, acc telegraf.Accumulator) {
	fields := map[string]interface{}{}
	tags := map[string]string{
		"resource":  r.Name,
		"namespace": r.Namespace,
	}

	for resourceName, val := range r.Status.Hard {
		switch resourceName {
		case "limits.cpu":
			fields["hard_cpu_cores_limit"] = ki.convertQuantity(val.String(), 1)
		case "limits.memory":
			fields["hard_memory_bytes_limit"] = ki.convertQuantity(val.String(), 1)
		case "requests.cpu":
			fields["hard_cpu_cores_request"] = ki.convertQuantity(val.String(), 1)
		case "requests.memory":
			fields["hard_memory_bytes_request"] = ki.convertQuantity(val.String(), 1)
		case "requests.storage":
			fields["hard_storage_bytes_request"] = ki.convertQuantity(val.String(), 1)
		}
	}

	for resourceName, val := range r.Status.Used {
		switch resourceName {
		case "limits.cpu":
			fields["used_cpu_cores_limit"] = ki.convertQuantity(val.String(), 1)
		case "limits.memory":
			fields["used_memory_bytes_limit"] = ki.convertQuantity(val.String(), 1)
		case "requests.cpu":
			fields["used_cpu_cores_request"] = ki.convertQuantity(val.String(), 1)
		case "requests.memory":
			fields["used_memory_bytes_request"] = ki.convertQuantity(val.String(), 1)
		case "requests.storage":
			fields["used_storage_bytes_request"] = ki.convertQuantity(val.String(), 1)
		}
	}

	acc.AddFields(resourcequotaMeasurement, fields, tags)
}
