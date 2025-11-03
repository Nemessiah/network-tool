package main

import (
	"net"
	"os"
	"strings"
	"testing"
)

func TestRenderSwitchIOS_Basic(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("10.1.0.0/24")
	if err != nil {
		t.Fatalf("bad test CIDR: %v", err)
	}

	var params = NetworkParams{
		Name:   "CLI_IT",
		VLANID: 1042,
		Subnet: ipNet,
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
		t.Fatalf("output mismatch:\n--- want ---\n%s\n--- got ---\n%s", want, got)
	}
}
