package main

import (
	"fmt"
	"net"
)

type NetworkParams struct {
	Name   string
	VLANID int
	Subnet *net.IPNet
}

func GenerateCommands(params NetworkParams) map[string]string {
	commands := make(map[string]string)

	commands["firewall"] = fmt.Sprintf(
		"set network %s subnet %s vlan %d",
		params.Name, params.Subnet, params.VLANID,
	)

	commands["switch"] = fmt.Sprintf(
		"vlan %d\nname %s",
		params.VLANID, params.Name,
	)

	commands["netbox"] = fmt.Sprintf(
		`{"name": "%s", "vid": %d, "prefix": "%s"}`,
		params.Name, params.VLANID, params.Subnet,
	)

	return commands
}
