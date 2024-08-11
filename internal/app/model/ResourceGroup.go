package model

type ResourceGroup struct {
	ResourceGroupName string            `json:"resourceGroupName"`
	Location          string            `json:"location"`
	Tags              map[string]string `json:"tags"`
}
