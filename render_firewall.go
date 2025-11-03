package main

import (
	"fmt"
)

func RenderFirewallPaloAlto(Interface NetworkParams) (string, error) {
	var err error

	output := fmt.Sprintf(
		"set network %s subnet %s vlan %d",
		Interface.Name, Interface.Subnet, Interface.VLANID,
	)

	return output, err
}
