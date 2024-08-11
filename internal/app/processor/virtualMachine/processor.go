package virtualMachine

import (
	"x/internal/app/clients"
	"x/internal/app/processor"
	"x/internal/app/processor/response"
	"x/internal/lib/api"
)

type Processor struct {
	clientRG     clients.ResourceGroup
	clientVNET   clients.VirtualNetwork
	clientSUBNET clients.Subnet
	clientIP     clients.PublicIPAddress
	clientNSG    clients.NetworkSecurityGroup
	clientNIC    clients.Interface
	clientVM     clients.VirtualMachine
}

func New(rg clients.ResourceGroup, vnet clients.VirtualNetwork,
	subnet clients.Subnet, ip clients.PublicIPAddress, nsg clients.NetworkSecurityGroup,
	nic clients.Interface, vm clients.VirtualMachine) *Processor {
	return &Processor{
		clientRG:     rg,
		clientVNET:   vnet,
		clientSUBNET: subnet,
		clientIP:     ip,
		clientNSG:    nsg,
		clientNIC:    nic,
		clientVM:     vm,
	}
}

func (p *Processor) Process(req *api.RequestAPI) response.ProcessedResponse {
	action := fetchAction(req.Action)
	switch action {
	case processor.CREATE:
		return p.processCreate(req.Data)
	case processor.UPDATE: //p.processUpdate(req.Data)
	case processor.DELETE: //p.processDelete(req.Data)
	case processor.LIST: //p.processList(req.Data
	case processor.UNKNOWN:
		//...
	}

	return nil
}

func fetchType(entityType string) processor.Type {
	switch entityType {
	case "resource group":
		return processor.ResourceGroup
	case "virtual network":
		return processor.Vnet
	case "subnet":
		return processor.Subnet
	case "public ip":
		return processor.PublicIP
	case "network security group":
		return processor.Nsg
	case "network interface":
		return processor.Nic
	case "virtual machine":
		return processor.Vm
	default:
		return processor.Unknown
	}
}

func fetchAction(action string) processor.Action {
	switch action {
	case "create":
		return processor.CREATE
	case "update":
		return processor.UPDATE
	case "delete":
		return processor.DELETE
	case "list":
		return processor.LIST
	default:
		return processor.UNKNOWN
	}
}
