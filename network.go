package main

import (
	"net"
)

type NetworkParams struct {
	Name   string
	VLANID int
	Subnet *net.IPNet
}

func GenerateCommands(params NetworkParams) (map[string]string, error) {
	commands := make(map[string]string)

	var err error

	commands["firewall"], err = RenderFirewallPaloAlto(params)

	commands["switch"], err = RenderSwitchIOS(params)

	commands["netbox"], err = Renderjson(params)

	return commands, err
}
