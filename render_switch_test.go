package main

import (
	"net"
	"os"
	"strings"
	"testing"
)

func visible(s string) string {
	s = strings.ReplaceAll(s, "\r", "␍")
	s = strings.ReplaceAll(s, "\n", "␤\n")
	s = strings.ReplaceAll(s, "\t", "→\t")
	return s
}
func TestRenderSwitchIOS_Basic(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("10.1.0.0/24")
	if err != nil {
		t.Fatalf("bad test CIDR: %v", err)
	}
	subnet := string(ipNet.IP)
	var params = NetworkParams{
		Name:   "CLI_IT",
		VLANID: 1042,
		Subnet: subnet,
	}

	got, err := RenderSwitchIOS(params)
	if err != nil {
		t.Fatalf("RenderSwitchIOS() returned error: %v", err)
	}

	wantBytes, err := os.ReadFile("testdata/switch_ios_basic.txt")
	want := string(wantBytes)

	got = strings.ReplaceAll(got, "\r\n", "\n")
	want = strings.ReplaceAll(want, "\r\n", "\n")

	if got != want {
		t.Fatalf("output mismatch:\n--- want ---\n%s\n--- got ---\n%s", visible(want), visible(got))
	}
}
