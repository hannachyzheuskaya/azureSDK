package virtualMachine

import (
	"context"
	"fmt"
	"sort"
	"x/internal/app/model"
	"x/internal/app/processor"
	"x/internal/app/processor/response"
	"x/internal/lib/api"
)

func (p *Processor) processCreate(entities []api.Entity) (resp response.ProcessedResponse) {
	var (
		resourceGroupName, location, vnetName           string
		subnetID, publicIPID, nsgID, networkInterfaceID string
		resourceResponse                                response.ProcessedResponse
	)

	sort.Slice(entities, func(i, j int) bool {
		type1 := fetchType(entities[i].Type)
		type2 := fetchType(entities[j].Type)
		return type1 < type2
	})

	for _, entity := range entities {
		entityType := fetchType(entity.Type)
		switch entityType {
		case processor.ResourceGroup:
			res, ok := entity.Meta.(model.ResourceGroup)
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion resource group error").
					Build()
			}
			resultRG, err := p.clientRG.Create(context.Background(), &res)
			if err != nil {
				fmt.Println(err)
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			resourceGroupName = *resultRG.Name
			location = *resultRG.Location
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultRG.ID).
				SetStatus(string(*resultRG.Properties.ProvisioningState)).
				Build()
		case processor.Vnet:
			res, ok := entity.Meta.(model.VirtualNetwork)
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion virtual network error").
					Build()
			}
			if res.ResourceGroupName == "" {
				res.ResourceGroupName = resourceGroupName
			} else {
				resourceGroupName = res.ResourceGroupName
			}
			if res.Location == "" {
				res.Location = location
			} else {
				location = res.Location
			}
			resultVN, err := p.clientVNET.Create(context.Background(), &res)
			if err != nil {
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			vnetName = *resultVN.Name
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultVN.ID).
				SetStatus(string(*resultVN.Properties.ProvisioningState)).
				Build()
		case processor.Subnet:
			res, ok := entity.Meta.(model.Subnet)
			res.VnetName = vnetName
			if res.ResourceGroupName == "" {
				res.ResourceGroupName = resourceGroupName
			} else {
				resourceGroupName = res.ResourceGroupName
			}
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion subnet error").
					Build()
			}
			resultSubnet, err := p.clientSUBNET.Create(context.Background(), &res)
			if err != nil {
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			subnetID = *resultSubnet.ID
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultSubnet.ID).
				SetStatus(string(*resultSubnet.Properties.ProvisioningState)).
				Build()
		case processor.PublicIP:
			res, ok := entity.Meta.(model.PublicIP)
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion public ip error").
					Build()
			}
			if res.ResourceGroupName == "" {
				res.ResourceGroupName = resourceGroupName
			} else {
				resourceGroupName = res.ResourceGroupName
			}
			if res.Location == "" {
				res.Location = location
			} else {
				location = res.Location
			}
			resultPublicIP, err := p.clientIP.Create(context.Background(), &res)
			if err != nil {
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			publicIPID = *resultPublicIP.ID
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultPublicIP.ID).
				SetStatus(string(*resultPublicIP.Properties.ProvisioningState)).
				Build()
		case processor.Nsg:
			res, ok := entity.Meta.(model.NetworkSecurityGroup)
			if res.ResourceGroupName == "" {
				res.ResourceGroupName = resourceGroupName
			} else {
				resourceGroupName = res.ResourceGroupName
			}
			if res.Location == "" {
				res.Location = location
			} else {
				location = res.Location
			}
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion network security group error").
					Build()
			}
			resultNsg, err := p.clientNSG.Create(context.Background(), &res)
			if err != nil {
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			nsgID = *resultNsg.ID
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultNsg.ID).
				SetStatus(string(*resultNsg.Properties.ProvisioningState)).
				Build()
		case processor.Nic:
			res, ok := entity.Meta.(model.NetworkInterface)
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion network interface error").
					Build()
			}
			if res.ResourceGroupName == "" {
				res.ResourceGroupName = resourceGroupName
			} else {
				resourceGroupName = res.ResourceGroupName
			}
			if res.Location == "" {
				res.Location = location
			} else {
				location = res.Location
			}
			res.SubnetID = subnetID
			res.PublicIpID = publicIPID
			res.NetworkSecurityGroupID = nsgID
			resultNIC, err := p.clientNIC.Create(context.Background(), &res)
			if err != nil {
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			networkInterfaceID = *resultNIC.ID
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultNIC.ID).
				SetStatus(string(*resultNIC.Properties.ProvisioningState)).
				Build()
		case processor.Vm:
			res, ok := entity.Meta.(model.VirtualMachine)
			if res.ResourceGroupName == "" {
				res.ResourceGroupName = resourceGroupName
			} else {
				resourceGroupName = res.ResourceGroupName
			}
			if res.Location == "" {
				res.Location = location
			} else {
				location = res.Location
			}
			res.NetworkInterfaceID = networkInterfaceID
			if !ok {
				return response.NewResponseBuilder().
					SetErrorDescription("conversion virtual machine error").
					Build()
			}
			resultVM, key, err := p.clientVM.Create(context.Background(), &res)
			if err != nil {
				return response.NewResponseBuilder().
					SetErrorDescription(err.Error()).
					Build()
			}
			resourceResponse = response.NewResponseBuilder().
				SetResourceID(*resultVM.ID).
				SetStatus(string(*resultVM.Properties.ProvisioningState)).
				SetPrivateKey(key).
				Build()

		case processor.Unknown:
			return response.NewResponseBuilder().
				SetErrorDescription("unknown operation entity").
				Build()
		}
	}
	return resourceResponse
}
