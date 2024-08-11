package model

type Subnet struct {
	ResourceGroupName string `json:"resourceGroupName"`
	VnetName          string `json:"vnetName"`
	SubnetName        string `json:"subnetName"`
	AddressPrefix     string `json:"addressPrefix"`
}
