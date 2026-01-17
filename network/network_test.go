package network

import (
	"strings"
	"testing"
)

func TestSelectInterface(t *testing.T) {
	tests := []struct {
		name        string
		networkType string
		vlan        int
		want        string
		wantErr     bool
		errString   string
	}{
		{
			name:        "WAN",
			networkType: "WAN",
			vlan:        1042,
			want:        "ae1.1042",
		},
		{
			name:        "Internal",
			networkType: "Internal",
			vlan:        1042,
			want:        "ae2.1042",
		},
		{
			name:        "Guest",
			networkType: "Guest",
			vlan:        1042,
			want:        "ae3.1042",
		},
		{
			name:        "DMZ",
			networkType: "DMZ",
			vlan:        1042,
			want:        "ae4.1042",
		},
		{
			name:        "default",
			networkType: "default",
			vlan:        1042,
			wantErr:     true,
			errString:   "invalid NetworkType: default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectInterface(tt.networkType, tt.vlan)

			if err == nil && tt.wantErr {
				if tt.errString != "" && !strings.Contains(err.Error(), tt.errString) {
					t.Fatalf("error = %q; want substring %q", err.Error(), tt.errString)
				}
			}
			if !strings.Contains(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestGetFirstUsableCIDR(t *testing.T) {
	tests := []struct {
		name      string
		cidr      string
		want      string
		wantErr   bool
		errString string
	}{
		{
			name: "Basic",
			cidr: "10.3.42.0/24",
			want: "10.3.42.1/24",
		},
		{
			name: "Basic_2",
			cidr: "10.3.42.10/29",
			want: "10.3.42.9/29",
		},
		{
			name:      "Invalid_CIDR",
			cidr:      "10.3.42.10_29",
			wantErr:   true,
			errString: "invalid CIDR",
		},
		{
			name:      "Invalid_NetworkSize",
			cidr:      "10.3.42.0/32",
			wantErr:   true,
			errString: "network too small, no usable host IPs",
		},
		{
			name:      "IPv6",
			cidr:      "fe80::88e8:5802:2d14:c001/64",
			wantErr:   true,
			errString: "only IPv4 is supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFirstUsableCIDR(tt.cidr)

			if err == nil && tt.wantErr {
				if tt.errString != "" && !strings.Contains(err.Error(), tt.errString) {
					t.Fatalf("error = %q; want substring %q", err.Error(), tt.errString)
				}
			}
			if !strings.Contains(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}
