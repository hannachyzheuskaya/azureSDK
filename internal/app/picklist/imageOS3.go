package picklist

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
)

type ImageOS3 struct {
	clientFactory *armcompute.ClientFactory
}

func NewImOS3(conn azcore.TokenCredential, subscriptionId string) (*ImageOS3, error) {
	cf, err := armcompute.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, errResourcesClientFactory
	}
	return &ImageOS3{
		clientFactory: cf,
	}, nil
}

func (i ImageOS3) GetImOS3(ctx context.Context, location string, publisherName string, offer string) ([]ImageAPI, error) {

	var imageOSs []ImageAPI
	res, err := i.clientFactory.NewVirtualMachineImagesClient().ListSKUs(ctx, location, publisherName, offer, nil)

	if err != nil {
		return nil, err
	}
	for _, v := range res.VirtualMachineImageResourceArray {
		imageOSs = append(imageOSs, convertIM3(v))
	}
	return imageOSs, nil

}

func convertIM3(from *armcompute.VirtualMachineImageResource) ImageAPI {
	return ImageAPI{
		Name: *from.Name,
	}
}
