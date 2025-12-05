package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/nemessiah/network-tool/internal"
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

func main() {
	var err error
	var visited map[string]bool

	configPath, err := internal.ConfigCheck()

	if err != nil {
		log.Fatalf("Unable to find, or create config: %v", err)
		return
	}

	Config, err := internal.LoadConfig(configPath, visited)

	if err := internal.ConfigValidation(&Config); err != nil {
		log.Fatalf("config validation failed: %v", err)
		return
	}

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

	Interface, err = network.SelectInterface(networkType, vlanID)

	if err != nil {
		fmt.Println("Error selecting interface:", err)
		return
	}

}
