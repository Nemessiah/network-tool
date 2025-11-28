package network

import (
	"fmt"
	"math/big"
	"net"
)

type NetworkParams struct {
	Name   string
	VLANID int
	Subnet string
}

func GetFirstUsableCIDR(cidr string) (string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR: %v", err)
	}

	networkIP := ipnet.IP.To4()
	if networkIP == nil {
		return "", fmt.Errorf("only IPv4 is supported")
	}

	// Convert network IP to big.Int
	ipInt := big.NewInt(0).SetBytes(networkIP)

	// Calculate total number of addresses
	maskSize, bits := ipnet.Mask.Size()
	totalIPs := 1 << (bits - maskSize)
	if totalIPs <= 2 {
		return "", fmt.Errorf("network too small, no usable host IPs")
	}

	// First usable IP = network + 1
	ipInt.Add(ipInt, big.NewInt(1))

	firstIP := net.IP(ipInt.Bytes())
	if len(firstIP) > 4 {
		firstIP = firstIP[len(firstIP)-4:]
	}

	// Build "[firstIP]/CIDRLength" string
	result := fmt.Sprintf("%s/%d", firstIP.String(), maskSize)
	return result, nil
}
