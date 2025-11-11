package firewall

import (
	"fmt"

	"github.com/nemessiah/network-tool/network"
)

func RenderFirewallPaloAlto(Interface network.NetworkParams) (string, error) {
	var err error

	output := fmt.Sprintf(
		"set network %s subnet %s vlan %d",
		Interface.Name, Interface.Subnet, Interface.VLANID,
	)

	return output, err
}
