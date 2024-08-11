package model

import (
	"encoding/json"
)

type VirtualMachine struct {
	ResourceGroupName  string         `json:"resource_group_name"`
	VmName             string         `json:"vm_name"`
	Location           string         `json:"location"`
	ImageOffer         string         `json:"image_offer"`
	ImagePublisher     string         `json:"image_publisher"`
	ImageSKU           string         `json:"image_sku"`
	ImageVersion       string         `json:"image_version"`
	VMSize             string         `json:"vm_size"`
	OSdiskName         string         `json:"osdisk_name"`
	NetworkInterfaceID string         `json:"network_interface_id"`
	Authentication     Authentication `json:"authentication"`
}

type Authentication struct {
	Type    string          `json:"authentication_type"`
	RawData json.RawMessage `json:"authentication_data"`
	Data    interface{}     `json:"-"`
}

type AuthenticationPSW struct {
	ComputerName  string `json:"computer_name"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
}

type AuthenticationSSH struct {
	AdminUsername    string `json:"admin_username"`
	ComputerName     string `json:"computer_name"`
	SshPublicKeyName string `json:"ssh_public_key_name"`
}

func (auth *Authentication) UnmarshalJSON(data []byte) error {
	type authentication Authentication
	var a authentication

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*auth = Authentication(a)
	switch a.Type {
	case "public key":
		var data AuthenticationSSH
		if err := json.Unmarshal(auth.RawData, &data); err != nil {
			return err
		}
		auth.Data = data
	case "password":
		var data AuthenticationPSW
		if err := json.Unmarshal(auth.RawData, &data); err != nil {
			return err
		}
		auth.Data = data
	}
	return nil
}
