package picklist

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
)

type ImageOS struct {
	clientFactory *armcompute.ClientFactory
}

func NewImOS(conn azcore.TokenCredential, subscriptionId string) (*ImageOS, error) {
	cf, err := armcompute.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, errResourcesClientFactory
	}
	return &ImageOS{
		clientFactory: cf,
	}, nil
}

func (i ImageOS) GetImOS(ctx context.Context, location string) ([]ImageAPI, error) {

	var imageOSs []ImageAPI
	res, err := i.clientFactory.NewVirtualMachineImagesClient().ListPublishers(ctx, location, nil)

	if err != nil {
		return nil, err
	}
	for _, v := range res.VirtualMachineImageResourceArray {
		imageOSs = append(imageOSs, convertIM(v))
	}
	return imageOSs, nil

}

func convertIM(from *armcompute.VirtualMachineImageResource) ImageAPI {
	return ImageAPI{
		Name: *from.Name,
	}
}
