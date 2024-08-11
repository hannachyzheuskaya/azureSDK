package model

type NetworkSecurityGroup struct {
	Name              string            `json:"name"`
	ResourceGroupName string            `json:"resourceGroupName"`
	Location          string            `json:"location"`
	Tags              map[string]string `json:"tags"`
	Rules             []SecurityRule    `json:"rules"`
}

type SecurityRule struct {
	Name                     string `json:"name"`
	SourceAddressPrefix      string `json:"source"`
	SourcePortRange          string `json:"sourcePortRange"`
	DestinationAddressPrefix string `json:"destination"`
	DestinationPortRange     string `json:"destinationPortRange"`
	Protocol                 string `json:"protocol"`
	Access                   string `json:"access"`
	Priority                 string `json:"priority"`
	Description              string `json:"description"`
	Direction                string `json:"direction"`
}
