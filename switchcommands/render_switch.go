package switchcommands

import (
	"fmt"

	"github.com/nemessiah/network-tool/network"
)

func RenderSwitchIOS(vlan network.NetworkParams) (string, error) {
	var err error

	output := fmt.Sprintf(
		"configure t\nvlan %d\nname %s",
		vlan.VLANID, vlan.Name,
	)

	return output, err
}
