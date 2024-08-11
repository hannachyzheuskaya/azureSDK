package publicIP

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/julien040/go-ternary"
	"x/internal/app/model"
)

type ClientPublicIP struct {
	client *armnetwork.PublicIPAddressesClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientPublicIP, error) {
	clientFactory, err := armnetwork.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, err
	}
	return &ClientPublicIP{
		client: clientFactory.NewPublicIPAddressesClient(),
	}, nil
}

func (r *ClientPublicIP) Create(ctx context.Context, ip *model.PublicIP) (*armnetwork.PublicIPAddress, error) {

	skuName := ternary.If(ip.SKU == "standard", armnetwork.PublicIPAddressSKUNameStandard, armnetwork.PublicIPAddressSKUNameBasic)
	publicIPAllocationMethod := ternary.If(ip.AllocationMethod == "Dynamic", armnetwork.IPAllocationMethodDynamic, armnetwork.IPAllocationMethodStatic)
	version := ternary.If(ip.Version == "IPv6", armnetwork.IPVersionIPv6, armnetwork.IPVersionIPv4)

	parameters := armnetwork.PublicIPAddress{
		Location: to.Ptr(ip.Location),
		SKU: &armnetwork.PublicIPAddressSKU{
			Name: to.Ptr(skuName),
		},
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: to.Ptr(publicIPAllocationMethod),
			PublicIPAddressVersion:   to.Ptr(version),
		},
	}

	pollerResponse, err := r.client.BeginCreateOrUpdate(ctx, ip.ResourceGroupName, ip.IPName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.PublicIPAddress, err
}

func (r *ClientPublicIP) Delete(ctx context.Context, ip *model.PublicIP) error {
	pollerResponse, err := r.client.BeginDelete(ctx, ip.ResourceGroupName, ip.IPName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientPublicIP) List(ctx context.Context, ip *model.PublicIP) ([]*armnetwork.PublicIPAddress, error) {
	//TODO implement
	return nil, nil
}
