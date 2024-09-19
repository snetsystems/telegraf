//go:generate ../../../tools/readme_config_includer/generator
// Package vrops implements an vrops input plugin for Telegraf
package vrops

import "regexp"

var (
	typeCloudZones         = regexp.MustCompile("^cloudzones:")
	typeCloudZonesTagValue = regexp.MustCompile("^cloudzones:(.*?)/Datacenter")
)

var (
	apiURL        = "/suite-api/api"
	authToken     = "/auth/token/acquire"
	resourceInfo  = "/resources/query"
	resourceBulk  = "/resources/bulk/relationships"
	resourceStats = "/resources/stats/latest/query"
)

var (
	vmStatKey = [...]string{"sys|osUptime_latest",
		"sys|poweredOn",
		"cpu|usage_average",
		"cpu|workload",
		"cpu|iowaitPct",
		"cpu|peak_vcpu_usage",
		"cpu|peak_vcpu_ready",
		"cpu|readyPct",
		"cpu|costopPct",
		"cpu|swapwaitPct",
		"cpu|vcpu_usage_disparity",
		"cpu|vm_capacity_provisioned",
		"cpu|effective_limit",
		"cpu|corecount_provisioned",
		"cpu|demandPct",
		"cpu|demandmhz",
		"cpu|capacity_contentionPct",
		"cpu|run_summation",
		"cpu|overlap_summation",
		"config|hardware|num_Cpu",
		"config|hardware|disk_Space",
		"mem-host|workload",
		"mem|usage_average",
		"mem|workload",
		"mem|guest_demand",
		"mem|swapped_average",
		"mem|swapoutRate_average",
		"mem|swapinRate_average",
		"mem|compressed_average",
		"mem|guest_provisioned",
		"mem|overhead_average",
		"mem|reservation_used",
		"mem|effective_limit",
		"diskspace|provisionedSpace",
		"diskspace|used",
		"diskspace|snapshot",
		"diskspace|snapshot|used",
		"diskspace|notshared",
		"diskspace|activeNotShared",
		"diskspace|workload",
		"guest|cpu_queue",
		"guest|contextSwapRate_latest",
		"guest|disk_queue",
		"guest|used_memory",
		"guest|mem.free_latest",
		"guest|mem.physUsable_latest",
		"guest|mem.needed_latest",
		"guest|page.outRate_latest",
		"guest|page.inRate_latest",
		"guest|page.size_latest",
		"guest|swap.spaceRemaining_latest",
		"guest|tools_running_status",
		"guestfilesystem|capacity_total",
		"guestfilesystem|freespace_total",
		"guestfilesystem|usage_total",
		"guestfilesystem|percentage_total",
		"net|transmitted_average",
		"net|received_average",
		"net|usage_average",
		"net|droppedTx_summation",
		"net|broadcastTx_summation",
		"net|multicastTx_summation",
		"net|droppedTx_summation_sum",
		"net|multicastTx_summation_sum",
		"net|broadcastTx_summation_sum",
		"storage|totalReadLatency_average",
		"storage|totalWriteLatency_average",
		"virtualDisk|read_average",
		"virtualDisk|write_average",
		"virtualDisk|totalLatency",
		"virtualDisk|totalReadLatency_average",
		"virtualDisk|totalWriteLatency_average",
		"virtualDisk|numberReadAveraged_average",
		"virtualDisk|numberWriteAveraged_average",
		"virtualDisk|commandsAveraged_average",
		"OnlineCapacityAnalytics|capacityRemainingPercentage",
		"OnlineCapacityAnalytics|timeRemaining",
		"OnlineCapacityAnalytics|cpu|capacityRemaining",
		"OnlineCapacityAnalytics|cpu|recommendedSize",
		"OnlineCapacityAnalytics|cpu|timeRemaining",
		"OnlineCapacityAnalytics|mem|capacityRemaining",
		"OnlineCapacityAnalytics|mem|recommendedSize",
		"OnlineCapacityAnalytics|mem|timeRemaining",
		"OnlineCapacityAnalytics|diskspace|capacityRemaining",
		"OnlineCapacityAnalytics|diskspace|recommendedSize",
		"OnlineCapacityAnalytics|diskspace|timeRemaining",
	}
)

type tokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type vROpsToken struct {
	Token     string        `json:"token"`
	Validity  int64         `json:"validity"`
	ExpiresAt string        `json:"expiresAt"`
	Roles     []interface{} `json:"roles"`
}

type resourceBody struct {
	ResourceKind   []string `json:"resourceKind"`
	ResourceState  []string `json:"resourceState"`
	ResourceStatus []string `json:"resourceStatus"`
}

type statBody struct {
	ResourceID  []string `json:"resourceId"`
	StatKey     []string `json:"statKey"`
	CurrentOnly bool     `json:"currentOnly"`
	MaxSamples  int      `json:"maxSamples"`
}

type relationshipBody struct {
	RelationshipType string       `json:"relationshipType"`
	ResourceIds      []string     `json:"resourceIds"`
	ResourceQuery    resourceBody `json:"resourceQuery"`
}

// resources is the main structure associated with a collection instance.
type resources struct {
	PageInfo struct {
		TotalCount int `json:"totalCount,omitempty"`
		Page       int `json:"page,omitempty"`
		PageSize   int `json:"pageSize,omitempty"`
	} `json:"pageInfo,omitempty"`
	Links []struct {
		Href string `json:"href,omitempty"`
		Rel  string `json:"rel,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"links,omitempty"`
	ResourceList []struct {
		CreationTime int64 `json:"creationTime,omitempty"`
		ResourceKey  struct {
			Name                string `json:"name,omitempty"`
			AdapterKindKey      string `json:"adapterKindKey,omitempty"`
			ResourceKindKey     string `json:"resourceKindKey,omitempty"`
			ResourceIdentifiers []struct {
				IdentifierType struct {
					Name               string `json:"name,omitempty"`
					DataType           string `json:"dataType,omitempty"`
					IsPartOfUniqueness bool   `json:"isPartOfUniqueness,omitempty"`
				} `json:"identifierType,omitempty"`
				Value string `json:"value,omitempty"`
			} `json:"resourceIdentifiers,omitempty"`
		} `json:"resourceKey,omitempty"`
		ResourceStatusStates []struct {
			AdapterInstanceID string `json:"adapterInstanceId,omitempty"`
			ResourceStatus    string `json:"resourceStatus,omitempty"`
			ResourceState     string `json:"resourceState,omitempty"`
			StatusMessage     string `json:"statusMessage,omitempty"`
		} `json:"resourceStatusStates,omitempty"`
		ResourceHealth      string  `json:"resourceHealth,omitempty"`
		ResourceHealthValue float64 `json:"resourceHealthValue,omitempty"`
		DtEnabled           bool    `json:"dtEnabled,omitempty"`
		Badges              []struct {
			Type  string  `json:"type,omitempty"`
			Color string  `json:"color,omitempty"`
			Score float64 `json:"score,omitempty"`
		} `json:"badges,omitempty"`
		RelatedResources []interface{} `json:"relatedResources,omitempty"`
		Links            []struct {
			Href string `json:"href,omitempty"`
			Rel  string `json:"rel,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"links,omitempty"`
		Identifier string `json:"identifier,omitempty"`
	} `json:"resourceList,omitempty"`
}

type bulkResource struct {
	PageInfo struct {
		TotalCount int `json:"totalCount"`
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
	} `json:"pageInfo"`
	RelationshipType   string `json:"relationshipType"`
	ResourcesRelations []struct {
		Resource struct {
			CreationTime int64 `json:"creationTime"`
			ResourceKey  struct {
				Name                string `json:"name"`
				AdapterKindKey      string `json:"adapterKindKey"`
				ResourceKindKey     string `json:"resourceKindKey"`
				ResourceIdentifiers []struct {
					IdentifierType struct {
						Name               string `json:"name"`
						DataType           string `json:"dataType"`
						IsPartOfUniqueness bool   `json:"isPartOfUniqueness"`
					} `json:"identifierType"`
					Value string `json:"value"`
				} `json:"resourceIdentifiers"`
			} `json:"resourceKey"`
			ResourceStatusStates []struct {
				AdapterInstanceID string `json:"adapterInstanceId"`
				ResourceStatus    string `json:"resourceStatus"`
				ResourceState     string `json:"resourceState"`
				StatusMessage     string `json:"statusMessage"`
			} `json:"resourceStatusStates"`
			ResourceHealth      string        `json:"resourceHealth"`
			ResourceHealthValue float64       `json:"resourceHealthValue"`
			DtEnabled           bool          `json:"dtEnabled"`
			Badges              []interface{} `json:"badges"`
			RelatedResources    []string      `json:"relatedResources"`
			Identifier          string        `json:"identifier"`
		} `json:"resource"`
		RelatedResources []string `json:"relatedResources"`
	} `json:"resourcesRelations"`
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
		Name string `json:"name"`
	} `json:"links"`
}

type statResource struct {
	Values []struct {
		ResourceID string `json:"resourceId"`
		StatList   struct {
			Stat []struct {
				Timestamps []int64 `json:"timestamps"`
				StatKey    struct {
					Key string `json:"key"`
				} `json:"statKey"`
				Data   []interface{} `json:"data,omitempty"`
				Values []string      `json:"values,omitempty"`
			} `json:"stat"`
		} `json:"stat-list"`
	} `json:"values"`
}
