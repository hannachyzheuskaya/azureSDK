package model

type NetworkInterface struct {
	ResourceGroupName         string            `json:"resourceGroupName"`
	Name                      string            `json:"name"`
	SubnetID                  string            `json:"subnetId"`
	PublicIpID                string            `json:"publicIp"`
	AllocationMethodPrivateIP string            `json:"allocationMethodPrivateIp"`
	IPName                    string            `json:"ipName"`
	Location                  string            `json:"location"`
	Tags                      map[string]string `json:"tags"`
	IPConfigurationName       string            `json:"ipConfigurationName"`
	NetworkSecurityGroupID    string            `json:"networkSecurityGroupId"`
}
