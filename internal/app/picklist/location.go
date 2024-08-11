package picklist

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type Location struct {
	clientFactory *armsubscriptions.ClientFactory
}

func NewL(conn azcore.TokenCredential) (*Location, error) {
	cf, err := armsubscriptions.NewClientFactory(conn, nil)
	if err != nil {
		return nil, errResourcesClientFactory
	}
	return &Location{
		clientFactory: cf,
	}, nil
}

func (l Location) GetL(ctx context.Context, subscriptionId string) ([]LocationAPI, error) {

	var locations []LocationAPI
	pager := l.clientFactory.NewClient().NewListLocationsPager(
		subscriptionId, &armsubscriptions.ClientListLocationsOptions{IncludeExtendedLocations: nil})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Value {
			if len(v.AvailabilityZoneMappings) != 0 {
				locations = append(locations, convertL(v))
			}
		}
	}
	return locations, nil
}

func convertL(from *armsubscriptions.Location) LocationAPI {
	return LocationAPI{
		Name:        *from.Name,
		DisplayName: *from.RegionalDisplayName,
	}

}
