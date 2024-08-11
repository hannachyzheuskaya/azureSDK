package model

type PublicIP struct {
	ResourceGroupName string            `json:"resourceGroupName"`
	AllocationMethod  string            `json:"allocationMethod"`
	IPName            string            `json:"ipName"`
	Location          string            `json:"location"`
	SKU               string            `json:"sku"`
	Version           string            `json:"version"`
	Tags              map[string]string `json:"tags"`
}
