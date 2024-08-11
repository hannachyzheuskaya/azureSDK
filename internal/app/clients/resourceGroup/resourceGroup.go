package resourceGroup

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"sync"
	"x/internal/app/model"
)

var instance *ClientResourceGroup
var once sync.Once
var subscription string
var credentials *azidentity.ClientSecretCredential

type ClientResourceGroup struct {
	client *armresources.ResourceGroupsClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientResourceGroup, error) {
	subscription = subscriptionId
	credentials = conn
	//singleton
	return getInstance()
}

func getInstance() (*ClientResourceGroup, error) {
	var err error
	var clientFactory *armresources.ClientFactory
	once.Do(func() {
		fmt.Println("creating new resource group client")
		clientFactory, err = armresources.NewClientFactory(subscription, credentials, nil)
		instance = &ClientResourceGroup{
			client: clientFactory.NewResourceGroupsClient(),
		}
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

/*
func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientResourceGroup, error) {
	fmt.Println("creating new resource group client")
	clientFactory, err := armresources.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, err
	}
	return &ClientResourceGroup{
		client: clientFactory.NewResourceGroupsClient(),
	}, nil
}
*/

func (r *ClientResourceGroup) Create(ctx context.Context, rg *model.ResourceGroup) (*armresources.ResourceGroup, error) {
	tagsPtrMap := make(map[string]*string)
	for k, v := range rg.Tags {
		tagsPtrMap[k] = to.Ptr(v)
	}

	parameters := armresources.ResourceGroup{
		Location: to.Ptr(rg.Location),
		Tags:     tagsPtrMap,
	}

	resp, err := r.client.CreateOrUpdate(ctx, rg.ResourceGroupName, parameters, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ResourceGroup, nil
}

func (r *ClientResourceGroup) Delete(ctx context.Context, rg *model.ResourceGroup) error {
	pollerResponse, err := r.client.BeginDelete(ctx, rg.ResourceGroupName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientResourceGroup) List(ctx context.Context, rg *model.ResourceGroup) ([]*armresources.ResourceGroup, error) {
	pager := r.client.NewListPager(nil)

	var resourceGroups []*armresources.ResourceGroup
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if nextResult.ResourceGroupListResult.Value != nil {
			resourceGroups = append(resourceGroups, nextResult.ResourceGroupListResult.Value...)
		}
	}

	return resourceGroups, nil
}
