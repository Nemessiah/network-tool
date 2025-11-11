package frontend

import (
	"github.com/nemessiah/network-tool/docs"
	"github.com/nemessiah/network-tool/firewall"
	"github.com/nemessiah/network-tool/network"
	"github.com/nemessiah/network-tool/switchcommands"
)

func GenerateCommands(params network.NetworkParams) (map[string]string, error) {
	commands := make(map[string]string)

	var err error

	commands["firewall"], err = firewall.RenderFirewallPaloAlto(params)

	commands["switch"], err = switchcommands.RenderSwitchIOS(params)

	commands["netbox"], err = docs.Renderjson(params)

	return commands, err
}
