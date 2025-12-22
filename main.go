package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/nemessiah/network-tool/commands"
	"github.com/nemessiah/network-tool/input/interactive"
	"github.com/nemessiah/network-tool/internal"
	"github.com/nemessiah/network-tool/network"
)

func main() {
	var err error
	var visited map[string]bool

	configPath, err := internal.ConfigCheck()

	if err != nil {
		log.Fatalf("Unable to find, or create config: %v", err)
		return
	}

	Config, err := internal.LoadConfig(configPath, visited)

	if err != nil {
		log.Fatalf("config load failed: %v", err)
		return
	}

	err = internal.ConfigValidation(&Config)

	if err != nil {
		log.Fatalf("config validation failed: %v", err)
		return
	}

	name := interactive.Prompt[string]("Enter VLAN Name: ")

	// VLAN ID with validation
	var vlanID int
	for {
		vlanID = interactive.Prompt[int]("Enter VLAN ID (1-4094): ")
		if vlanID < 1 || vlanID > 4094 {
			fmt.Println("Invalid VLAN ID, must be 1-4094.")
			continue
		}
		break
	}

	// Subnet input with CIDR validation
	var subnet string
	for {
		subnetStr := interactive.Prompt[string]("Enter Subnet (e.g. 10.1.0.0/24): ")
		_, netCIDR, err := net.ParseCIDR(subnetStr)
		if err != nil {
			fmt.Printf("Invalid subnet (%s) format. Please use CIDR notation (e.g., 10.1.0.0/24).\n", subnetStr)
			continue
		}
		subnet = netCIDR.String()
		break
	}

	var ipaddress string
	ipaddress, err = network.GetFirstUsableCIDR(subnet)

	if err != nil {
		fmt.Printf("Error locating first usable ip: %s", err)
	}

	// Network type validation
	var networkType string
	validTypes := map[string]bool{"WAN": true, "Internal": true, "Guest": true, "DMZ": true}
	for {
		networkType = interactive.Prompt[string]("Enter network type (WAN, Internal, Guest, DMZ): ")
		if !validTypes[networkType] {
			fmt.Println("Zone must be one of WAN, Internal, Guest, DMZ")
			continue
		}
		break
	}

	// Network description validation (max 50 characters)
	var description string
	for {
		description = interactive.Prompt[string]("Enter network description (50 Char max): ")
		if len(description) > 50 {
			fmt.Println("Max description length is 50 characters.")
			continue
		}
		break
	}

	// Zone input validation (max 50 characters)
	var zone string
	for {
		zone = interactive.Prompt[string]("Enter zone: ")
		if len(zone) > 50 {
			fmt.Println("Max zone length is 50 characters.")
			continue
		}
		break
	}

	// interface selection, and validation
	var Interface string

	Interface, err = network.SelectInterface(networkType, vlanID)

	if err != nil {
		fmt.Println("Error selecting interface:", err)
		return
	}

	// action input and validation
	var action string
	validCrud := map[string]struct{}{
		"create": {},
		"read":   {},
		"update": {},
		"delete": {},
	}

	for {
		action = interactive.Prompt[string]("Enter wanted CRUD action: ")
		action = strings.ToLower(strings.TrimSpace(action))
		if _, ok := validCrud[action]; ok {
			break
		}
		fmt.Println("Invalid action. Please try again.")
	}

	var extractedCommands map[string][]string
	extractedCommands = commands.ExtractVendorCommandsForAction(Config, action)

	network := commands.Network{
		VlanName: name,
		Vlanid:   vlanID,
		Subnet:   subnet,
	}

	inputs := commands.NetworkParams{
		Interface:   Interface,
		Network:     network,
		NetworkType: networkType,
		Description: description,
		Zone:        zone,
		IpAddress:   ipaddress,
	}

	outputCommands := make(map[string][]string)

	reflectedinput, err := commands.MakeStructMap(inputs)
	if err != nil {
		fmt.Printf("Error refecting input: %s", err)
	}

	for vendor, commandArray := range extractedCommands {
		replacedCommands := make([]string, 0, len(commandArray))
		for _, command := range commandArray {
			temp, err := commands.ReplaceKeys(reflectedinput, command)

			if err != nil {
				fmt.Printf("Error replacing keys in a command: %s", err)
			}
			replacedCommands = append(replacedCommands, temp)
		}
		outputCommands[vendor] = replacedCommands
	}

	for vendor, commands := range outputCommands {
		fmt.Println("commands for:", vendor)
		for _, command := range commands {
			fmt.Println("    ", command)
		}
		fmt.Println()
	}

}
