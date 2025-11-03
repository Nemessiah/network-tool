package main

import (
	"fmt"
)

func RenderSwitchIOS(vlan NetworkParams) (string, error) {
	var err error

	output := fmt.Sprintf(
		"configure t\nvlan %d\nname %s",
		vlan.VLANID, vlan.Name,
	)

	return output, err
}
