package main

import (
	"bufio" // for reading text input
	"fmt"   // for printing and formatting
	"net"
	"os"      // gives access to stdin/stdout
	"strconv" // convert strings to numbers
	"strings" // string trimming/splitting

	"github.com/nemessiah/network-tool/firewall"
	"github.com/nemessiah/network-tool/frontend"
	"github.com/nemessiah/network-tool/network"
)

func Prompt[T any](prompt string) T {
	reader := bufio.NewReader(os.Stdin)
	var zero T // default zero value of type T

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch any(zero).(type) {
		case int:
			val, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid entry, must be an integer.")
				continue
			}
			return any(val).(T)
		case string:
			return any(input).(T)
		default:
			panic("unsupported type")
		}
	}
}

func SelectInterface(Type string, Vlan int) (string, error) {
	var output string
	var err error

	switch Type {
	case "WAN":
		output = fmt.Sprintf("ae1.%d", Vlan)
	case "Internal":
		output = fmt.Sprintf("ae2.%d", Vlan)
	case "Guest":
		output = fmt.Sprintf("ae3.%d", Vlan)
	case "DMZ":
		output = fmt.Sprintf("ae4.%d", Vlan)
	default:
		return "", fmt.Errorf("invalid NetworkType: %s", Type)
	}

	return output, err
}

func main() {

	name := Prompt[string]("Enter VLAN Name: ")

	// VLAN ID with validation
	var vlanID int
	for {
		vlanID = Prompt[int]("Enter VLAN ID (1-4094): ")
		if vlanID < 1 || vlanID > 4094 {
			fmt.Println("Invalid VLAN ID, must be 1-4094.")
			continue
		}
		break
	}

	// Subnet input with CIDR validation
	var subnet string
	for {
		subnetStr := Prompt[string]("Enter Subnet (e.g. 10.1.0.0/24): ")
		_, netCIDR, err := net.ParseCIDR(subnetStr)
		if err != nil {
			fmt.Printf("Invalid subnet (%s) format. Please use CIDR notation (e.g., 10.1.0.0/24).\n", subnetStr)
			continue
		}
		subnet = netCIDR.String()
		break
	}

	// Network type validation
	var networkType string
	validTypes := map[string]bool{"WAN": true, "Internal": true, "Guest": true, "DMZ": true}
	for {
		networkType = Prompt[string]("Enter network type (WAN, Internal, Guest, DMZ): ")
		if !validTypes[networkType] {
			fmt.Println("Zone must be one of WAN, Internal, Guest, DMZ")
			continue
		}
		break
	}

	// Network description validation (max 50 characters)
	var description string
	for {
		description = Prompt[string]("Enter network description: ")
		if len(description) > 50 {
			fmt.Println("Max description length is 50 characters.")
			continue
		}
		break
	}

	// Zone input validation (max 50 characters)
	var zone string
	for {
		zone = Prompt[string]("Enter zone: ")
		if len(zone) > 50 {
			fmt.Println("Max zone length is 50 characters.")
			continue
		}
		break
	}

	var Interface string
	var err error

	Interface, err = SelectInterface(networkType, vlanID)

	if err != nil {
		fmt.Println("Error selecting interface:", err)
		return
	}

	var networkParams network.NetworkParams

	networkParams = network.NetworkParams{
		Name:   name,
		VLANID: vlanID,
		Subnet: subnet,
	}

	var firewallParams firewall.InterfaceConfig

	firewallParams = firewall.InterfaceConfig{
		Interface:   Interface,
		Network:     networkParams,
		NetworkType: networkType,
		Description: description,
		Zone:        zone,
	}

	commands, err := frontend.GenerateCommands(firewallParams)

	if err != nil {
		fmt.Println("Error generating commands:", err)
		return
	}

	fmt.Println("\n--- Generated Commands ---")
	fmt.Printf("Firewall:\n%s\n\n", commands["firewall"])
	fmt.Printf("Switch:\n%s\n\n", commands["switch"])
	fmt.Printf("NetBox JSON:\n%s\n", commands["netbox"])
}
