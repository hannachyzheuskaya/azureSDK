package virtualMachine

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"x/internal/app/model"
)

type ClientVirtualMachine struct {
	client     *armcompute.VirtualMachinesClient
	clientKeys *armcompute.SSHPublicKeysClient
	//disksClient *armcompute.DisksClient
}

func New(subscriptionId string, conn *azidentity.ClientSecretCredential) (*ClientVirtualMachine, error) {
	clientFactory, err := armcompute.NewClientFactory(subscriptionId, conn, nil)

	if err != nil {
		return nil, err
	}
	return &ClientVirtualMachine{
		client:     clientFactory.NewVirtualMachinesClient(),
		clientKeys: clientFactory.NewSSHPublicKeysClient(),
		//disksClient : clientFactory.NewDisksClient(),
	}, nil
}

func (r *ClientVirtualMachine) Create(ctx context.Context, machine *model.VirtualMachine) (*armcompute.VirtualMachine, string, error) {

	var privateKey string
	var authenticationMethod *armcompute.OSProfile

	if machine.Authentication.Type == "password" {
		auth, ok := machine.Authentication.Data.(model.AuthenticationPSW)
		if !ok { //TODO
			fmt.Println("TODO: handle error")
		}
		authenticationMethod = &armcompute.OSProfile{
			ComputerName:  to.Ptr(auth.ComputerName),
			AdminUsername: to.Ptr(auth.AdminUsername),
			AdminPassword: to.Ptr(auth.AdminPassword),
		}
	} else if machine.Authentication.Type == "public key" {
		auth, ok := machine.Authentication.Data.(model.AuthenticationSSH)
		r.clientKeys.Create(ctx, machine.ResourceGroupName, auth.SshPublicKeyName, armcompute.SSHPublicKeyResource{
			Location: to.Ptr(machine.Location),
		}, nil)
		res, err := r.clientKeys.GenerateKeyPair(ctx, machine.ResourceGroupName, auth.SshPublicKeyName, &armcompute.SSHPublicKeysClientGenerateKeyPairOptions{Parameters: nil})
		if !ok { //TODO
			fmt.Println("TODO: handle error")
		}
		if err != nil {
			//TODO
			fmt.Println("TODO: handle error")
		}
		privateKey = *res.PrivateKey
		authenticationMethod = &armcompute.OSProfile{
			ComputerName:  to.Ptr(auth.ComputerName),
			AdminUsername: to.Ptr(auth.AdminUsername),
			LinuxConfiguration: &armcompute.LinuxConfiguration{
				DisablePasswordAuthentication: to.Ptr(true),
				SSH: &armcompute.SSHConfiguration{
					PublicKeys: []*armcompute.SSHPublicKey{
						{
							Path:    to.Ptr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", auth.AdminUsername)),
							KeyData: res.PublicKey,
						},
					},
				},
			},
		}
	}

	parameters := armcompute.VirtualMachine{
		Location: to.Ptr(machine.Location),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Offer:     to.Ptr(machine.ImageOffer),
					Publisher: to.Ptr(machine.ImagePublisher),
					SKU:       to.Ptr(machine.ImageSKU),
					Version:   to.Ptr(machine.ImageVersion),
				},

				//TODO: hardcoded values
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(machine.OSdiskName),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
					},
					//DiskSizeGB: to.Ptr[int32](100), // default 127G
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes(machine.VMSize)), // VM size include vCPUs,RAM,Data Disks,Temp storage.
			},
			OSProfile: authenticationMethod,
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: to.Ptr(machine.NetworkInterfaceID),
					},
				},
			},
		},
	}

	pollerResponse, err := r.client.BeginCreateOrUpdate(ctx, machine.ResourceGroupName, machine.VmName, parameters, nil)
	if err != nil {
		return nil, "-", err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, "-", err
	}

	return &resp.VirtualMachine, privateKey, nil
}

func (r *ClientVirtualMachine) Delete(ctx context.Context, machine *model.VirtualMachine) error {
	pollerResponse, err := r.client.BeginDelete(ctx, machine.ResourceGroupName, machine.VmName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientVirtualMachine) List(ctx context.Context, ip *model.VirtualMachine) ([]*armcompute.VirtualMachine, error) {
	//TODO implement
	return nil, nil
}
