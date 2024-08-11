package picklist

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
)

type Size struct {
	clientFactory *armcompute.ClientFactory
}

func NewS(conn azcore.TokenCredential, subscriptionId string) (*Size, error) {
	cf, err := armcompute.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, errResourcesClientFactory
	}
	return &Size{
		clientFactory: cf,
	}, nil
}

func (s Size) GetS(ctx context.Context, location string) ([]SizeAPI, error) {

	var sizes []SizeAPI
	pager := s.clientFactory.NewResourceSKUsClient().NewListPager(&armcompute.ResourceSKUsClientListOptions{Filter: to.Ptr(location),
		IncludeExtendedLocations: nil,
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Value {
			sizes = append(sizes, convertS(v))
		}
	}
	return sizes, nil
}

func convertS(fromS *armcompute.ResourceSKU) SizeAPI {
	return SizeAPI{
		Name: *fromS.Name,
	}
}
