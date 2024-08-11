package subnet

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"

	"x/internal/app/model"
)

type ClientSubnet struct {
	client *armnetwork.SubnetsClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientSubnet, error) {
	clientFactory, err := armnetwork.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, err
	}
	return &ClientSubnet{
		client: clientFactory.NewSubnetsClient(),
	}, nil
}

func (r *ClientSubnet) Create(ctx context.Context, snet *model.Subnet) (*armnetwork.Subnet, error) {
	parameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.Ptr(snet.AddressPrefix),
		},
	}

	pollerResponse, err := r.client.BeginCreateOrUpdate(ctx, snet.ResourceGroupName, snet.VnetName, snet.SubnetName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Subnet, nil
}

func (r *ClientSubnet) Delete(ctx context.Context, snet *model.Subnet) error {
	pollerResponse, err := r.client.BeginDelete(ctx, snet.ResourceGroupName, snet.VnetName, snet.SubnetName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientSubnet) List(ctx context.Context, snet *model.Subnet) ([]*armnetwork.Subnet, error) {
	pager := r.client.NewListPager(snet.ResourceGroupName, snet.VnetName, nil)

	var vn []*armnetwork.Subnet
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if nextResult.SubnetListResult.Value != nil {
			vn = append(vn, nextResult.SubnetListResult.Value...)
		}
	}

	return vn, nil
}
