package docs

import (
	"encoding/json"

	"github.com/nemessiah/network-tool/network"
)

type NetboxPayload struct {
	Name   string `json:"name"`
	VID    int    `json:"vid"`
	Prefix string `json:"prefix"`
}

func Renderjson(info network.NetworkParams) (string, error) {
	var err error

	jsonInput := NetboxPayload{
		Name:   info.Name,
		VID:    info.VLANID,
		Prefix: info.Subnet,
	}

	json, err := json.Marshal(jsonInput)

	output := string(json)

	return output, err
}
