package frontend

import (
	"github.com/nemessiah/network-tool/docs"
	"github.com/nemessiah/network-tool/firewall"
	"github.com/nemessiah/network-tool/switchcommands"
)

func GenerateCommands(params firewall.InterfaceConfig) (map[string]string, error) {
	commands := make(map[string]string)

	var err error

	commands["firewall"], err = firewall.RenderFirewallPaloAlto(params)

	commands["switch"], err = switchcommands.RenderSwitchIOS(params.Network)

	commands["netbox"], err = docs.Renderjson(params.Network)

	return commands, err
}
