package picklist

import "errors"

var (
	errResourcesClientFactory = errors.New("failed to create resources client factory")
)

type LocationAPI struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type SizeAPI struct {
	Name string `json:"name"`
}

type ImageAPI struct {
	Name string `json:"name"`
}
