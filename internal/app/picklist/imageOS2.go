package picklist

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
)

type ImageOS2 struct {
	clientFactory *armcompute.ClientFactory
}

func NewImOS2(conn azcore.TokenCredential, subscriptionId string) (*ImageOS2, error) {
	cf, err := armcompute.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, errResourcesClientFactory
	}
	return &ImageOS2{
		clientFactory: cf,
	}, nil
}

func (i ImageOS2) GetImOS2(ctx context.Context, location string, publisherName string) ([]ImageAPI, error) {

	var imageOSs []ImageAPI
	res, err := i.clientFactory.NewVirtualMachineImagesClient().ListOffers(ctx, location, publisherName, nil)

	if err != nil {
		return nil, err
	}
	for _, v := range res.VirtualMachineImageResourceArray {
		imageOSs = append(imageOSs, convertIM2(v))
	}
	return imageOSs, nil

}

func convertIM2(from *armcompute.VirtualMachineImageResource) ImageAPI {
	return ImageAPI{
		Name: *from.Name,
	}
}
