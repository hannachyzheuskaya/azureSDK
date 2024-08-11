package virtualNetwork

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"x/internal/app/model"
)

type ClientVirtualNetwork struct {
	client *armnetwork.VirtualNetworksClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientVirtualNetwork, error) {
	clientFactory, err := armnetwork.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, err
	}
	return &ClientVirtualNetwork{
		client: clientFactory.NewVirtualNetworksClient(),
	}, nil
}

func (r *ClientVirtualNetwork) Create(ctx context.Context, vn *model.VirtualNetwork) (*armnetwork.VirtualNetwork, error) {
	addressPrefixes := make([]*string, len(vn.AddressSpace))
	for i, ap := range vn.AddressSpace {
		addressPrefixes[i] = to.Ptr(ap)
	}

	parameters := armnetwork.VirtualNetwork{
		Location: to.Ptr(vn.Location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: addressPrefixes,
			},
		},
	}

	pollerResponse, err := r.client.BeginCreateOrUpdate(ctx, vn.ResourceGroupName, vn.VnetName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualNetwork, nil
}

func (r *ClientVirtualNetwork) Delete(ctx context.Context, vn *model.VirtualNetwork) error {
	pollerResponse, err := r.client.BeginDelete(ctx, vn.ResourceGroupName, vn.VnetName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientVirtualNetwork) List(ctx context.Context, vnet *model.VirtualNetwork) ([]*armnetwork.VirtualNetwork, error) {
	pager := r.client.NewListPager(vnet.ResourceGroupName, nil)

	var vn []*armnetwork.VirtualNetwork
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if nextResult.VirtualNetworkListResult.Value != nil {
			vn = append(vn, nextResult.VirtualNetworkListResult.Value...)
		}
	}

	return vn, nil
}
