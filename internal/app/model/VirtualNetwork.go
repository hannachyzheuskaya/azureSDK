package model

type VirtualNetwork struct {
	ResourceGroupName string            `json:"resourceGroupName"`
	VnetName          string            `json:"vnetName"`
	Location          string            `json:"location"`
	AddressSpace      []string          `json:"addressSpace"`
	Tags              map[string]string `json:"tags"`
}
