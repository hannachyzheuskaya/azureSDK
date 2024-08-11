package networkSecurityGroup

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"x/internal/app/model"
)

type ClientSecurityGroup struct {
	client *armnetwork.SecurityGroupsClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientSecurityGroup, error) {
	clientFactory, err := armnetwork.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return nil, err
	}
	return &ClientSecurityGroup{
		client: clientFactory.NewSecurityGroupsClient(),
	}, nil
}

func (r *ClientSecurityGroup) Create(ctx context.Context, nsg *model.NetworkSecurityGroup) (*armnetwork.SecurityGroup, error) {

	rules := make([]*armnetwork.SecurityRule, 0, len(nsg.Rules))
	for _, rule := range nsg.Rules {
		rules = append(rules, convert(rule))
	}

	parameters := armnetwork.SecurityGroup{
		Location: to.Ptr(nsg.Location),
		Properties: &armnetwork.SecurityGroupPropertiesFormat{
			SecurityRules: rules,
		},
	}

	pollerResponse, err := r.client.BeginCreateOrUpdate(ctx, nsg.ResourceGroupName, nsg.Name, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.SecurityGroup, nil
}

func (r *ClientSecurityGroup) Delete(ctx context.Context, nsg *model.NetworkSecurityGroup) error {

	pollerResponse, err := r.client.BeginDelete(ctx, nsg.ResourceGroupName, nsg.Name, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClientSecurityGroup) List(ctx context.Context, nsg *model.NetworkSecurityGroup) ([]*armnetwork.SecurityGroup, error) {
	pager := r.client.NewListPager(nsg.ResourceGroupName, nil)

	var securityGroups []*armnetwork.SecurityGroup
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if nextResult.SecurityGroupListResult.Value != nil {
			securityGroups = append(securityGroups, nextResult.SecurityGroupListResult.Value...)
		}
	}

	return securityGroups, nil
}
