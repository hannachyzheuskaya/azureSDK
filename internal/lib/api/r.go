package api

import (
	"encoding/json"
	"x/internal/app/model"
)

type RequestAPI struct {
	Action string   `json:"action"`
	Data   []Entity `json:"entities"`
}

type Entity struct {
	Type    string          `json:"type"`
	RawData json.RawMessage `json:"meta"`
	Meta    interface{}     `json:"-"`
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type entity Entity
	var en entity

	if err := json.Unmarshal(data, &en); err != nil {
		return err
	}
	*e = Entity(en)
	switch e.Type {
	case "resource group":
		var data model.ResourceGroup
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	case "virtual network":
		var data model.VirtualNetwork
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	case "subnet":
		var data model.Subnet
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	case "public ip":
		var data model.PublicIP
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	case "network security group":
		var data model.NetworkSecurityGroup
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	case "network interface":
		var data model.NetworkInterface
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	case "virtual machine":
		var data model.VirtualMachine
		if err := json.Unmarshal(e.RawData, &data); err != nil {
			return err
		}
		e.Meta = data
	}
	return nil
}
