package clients

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"x/internal/app/model"
)

type ResourceGroup interface {
	Create(ctx context.Context, resource *model.ResourceGroup) (*armresources.ResourceGroup, error)
	Delete(ctx context.Context, resource *model.ResourceGroup) error
	List(ctx context.Context, resource *model.ResourceGroup) ([]*armresources.ResourceGroup, error)
}

type VirtualNetwork interface {
	Create(ctx context.Context, resource *model.VirtualNetwork) (*armnetwork.VirtualNetwork, error)
	Delete(ctx context.Context, resource *model.VirtualNetwork) error
	List(ctx context.Context, resource *model.VirtualNetwork) ([]*armnetwork.VirtualNetwork, error)
}

type Subnet interface {
	Create(ctx context.Context, resource *model.Subnet) (*armnetwork.Subnet, error)
	Delete(ctx context.Context, resource *model.Subnet) error
	List(ctx context.Context, resource *model.Subnet) ([]*armnetwork.Subnet, error)
}

type PublicIPAddress interface {
	Create(ctx context.Context, resource *model.PublicIP) (*armnetwork.PublicIPAddress, error)
	Delete(ctx context.Context, resource *model.PublicIP) error
	List(ctx context.Context, resource *model.PublicIP) ([]*armnetwork.PublicIPAddress, error)
}

type NetworkSecurityGroup interface {
	Create(ctx context.Context, resource *model.NetworkSecurityGroup) (*armnetwork.SecurityGroup, error)
	Delete(ctx context.Context, resource *model.NetworkSecurityGroup) error
	List(ctx context.Context, resource *model.NetworkSecurityGroup) ([]*armnetwork.SecurityGroup, error)
}

type Interface interface {
	Create(ctx context.Context, resource *model.NetworkInterface) (*armnetwork.Interface, error)
	Delete(ctx context.Context, resource *model.NetworkInterface) error
	List(ctx context.Context, resource *model.NetworkInterface) ([]*armnetwork.Interface, error)
}

type VirtualMachine interface {
	Create(ctx context.Context, machine *model.VirtualMachine) (*armcompute.VirtualMachine, string, error)
	Delete(ctx context.Context, resource *model.VirtualMachine) error
	List(ctx context.Context, resource *model.VirtualMachine) ([]*armcompute.VirtualMachine, error)
}
