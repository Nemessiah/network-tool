package main

import (
	"bufio" // for reading text input
	"fmt"   // for printing and formatting
	"net"
	"os"      // gives access to stdin/stdout
	"strconv" // convert strings to numbers
	"strings" // string trimming/splitting

	"github.com/nemessiah/network-tool/frontend"
	"github.com/nemessiah/network-tool/network"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter VLAN Name: ")
	name, _ := reader.ReadString('\n') // read until user presses Enter
	name = strings.TrimSpace(name)     // remove \n

	fmt.Print("Enter VLAN ID (1-4094): ")
	vlanStr, _ := reader.ReadString('\n')
	vlanStr = strings.TrimSpace(vlanStr)
	vlanID, err := strconv.Atoi(vlanStr) // convert to int

	if err != nil {
		fmt.Println("Invalid VLAN ID, must be a number.")
		return
	}

	if vlanID < 1 || vlanID > 4096 {
		fmt.Println("Invalid VLAN ID, must be a 1-4096.")
		return
	}

	fmt.Print("Enter Subnet (e.g. 10.1.0.0/24): ")
	subnetStr, _ := reader.ReadString('\n')
	subnetStr = strings.TrimSpace(subnetStr)

	_, net, err := net.ParseCIDR(subnetStr)
	if err != nil {
		fmt.Println("Invalid subnet (" + subnetStr + ") format. Please use CIDR notation (e.g., 10.1.0.0/24).")
		return
	}
	subnet := string(net.String())
	fmt.Println("Using network:", subnet)

	var params network.NetworkParams

	params = network.NetworkParams{
		Name:   name,
		VLANID: vlanID,
		Subnet: subnet,
	}

	commands, err := frontend.GenerateCommands(params)

	fmt.Println("\n--- Generated Commands ---")
	fmt.Printf("Firewall:\n%s\n\n", commands["firewall"])
	fmt.Printf("Switch:\n%s\n\n", commands["switch"])
	fmt.Printf("NetBox JSON:\n%s\n", commands["netbox"])
}
