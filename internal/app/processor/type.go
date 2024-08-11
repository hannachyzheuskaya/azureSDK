package processor

import (
	"x/internal/lib/api"
)

type Type int

const (
	Unknown Type = iota
	ResourceGroup
	Vnet
	Subnet
	PublicIP
	Nsg
	Nic
	Vm
)

type Action int

const (
	UNKNOWN Action = iota
	CREATE
	UPDATE
	DELETE
	LIST
)

type Processor interface {
	Process(req *api.RequestAPI) error
}
