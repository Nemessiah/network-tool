package firewall

import (
	"fmt"

	"github.com/nemessiah/network-tool/network"
)

// Intended output:
// configure
// set network interface aggregate-ethernet [AE_INTERFACE].[SUBINT_VLAN] layer3 ip [IP_ADDRESS/MASK] vlan [SUBINT_VLAN]
// set network interface aggregate-ethernet [AE_INTERFACE].[SUBINT_VLAN] description "[DESCRIPTION]"
// set zone [ZONE_NAME] network layer3 [AE_INTERFACE].[SUBINT_VLAN]
// commit

type InterfaceConfig struct {
	Interface   string                `json:"interface"`
	Network     network.NetworkParams `json:"network"`
	NetworkType string                `json:"network_type"`
	Description string                `json:"description"`
	Zone        string                `json:"zone"`
}

func RenderFirewallPaloAlto(Interface InterfaceConfig) (string, error) {
	var err error

	interfaceConfig, err := createPaloAltoInterface(Interface)

	descriptionConfig, err := SetPaloAltoDescription(Interface)

	zoneConfig, err := SetPaloAltoZone(Interface)

	output := fmt.Sprintf(
		"%s\n%s\n%s",
		interfaceConfig, descriptionConfig, zoneConfig,
	)
	return output, err
}

func createPaloAltoInterface(Interface InterfaceConfig) (string, error) {
	var err error

	firstUsable, err := network.GetFirstUsableCIDR(Interface.Network.Subnet)

	output := fmt.Sprintf(
		"set network interface aggregate-ethernet %s layer3 ip %s vlan %d",
		Interface.Interface, firstUsable, Interface.Network.VLANID,
	)

	return output, err
}

func SetPaloAltoDescription(Interface InterfaceConfig) (string, error) {
	var err error

	output := fmt.Sprintf(
		"set network interface aggregate-ethernet %s description %q",
		Interface.Interface, Interface.Description,
	)

	return output, err
}

func SetPaloAltoZone(Interface InterfaceConfig) (string, error) {
	var err error

	output := fmt.Sprintf(
		"set zone %s network layer3 %s",
		Interface.Zone, Interface.Interface,
	)

	return output, err
}
