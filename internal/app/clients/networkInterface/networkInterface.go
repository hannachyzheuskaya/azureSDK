package networkInterface

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"x/internal/app/model"
)

type ClientNetworkInterface struct {
	client *armnetwork.InterfacesClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientNetworkInterface, error) {
	clientFactory, err := armnetwork.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, err
	}
	return &ClientNetworkInterface{
		client: clientFactory.NewInterfacesClient(),
	}, nil
}

func (r *ClientNetworkInterface) Create(ctx context.Context, netInterface *model.NetworkInterface) (*armnetwork.Interface, error) {
	parameters := armnetwork.Interface{
		Location: to.Ptr(netInterface.Location),
		Properties: &armnetwork.InterfacePropertiesFormat{
			IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
				{
					Name: to.Ptr(netInterface.IPConfigurationName),
					Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						Subnet: &armnetwork.Subnet{
							ID: to.Ptr(netInterface.SubnetID),
						},
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID: to.Ptr(netInterface.PublicIpID),
						},
					},
				},
			},
			NetworkSecurityGroup: &armnetwork.SecurityGroup{
				ID: to.Ptr(netInterface.NetworkSecurityGroupID),
			},
		},
	}

	pollerResponse, err := r.client.BeginCreateOrUpdate(ctx, netInterface.ResourceGroupName, netInterface.Name, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Interface, err
}

func (r *ClientNetworkInterface) Delete(ctx context.Context, netInterface *model.NetworkInterface) error {
	pollerResponse, err := r.client.BeginDelete(ctx, netInterface.ResourceGroupName, netInterface.Name, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientNetworkInterface) List(ctx context.Context, ip *model.NetworkInterface) ([]*armnetwork.Interface, error) {
	//TODO implement
	return nil, nil
}
