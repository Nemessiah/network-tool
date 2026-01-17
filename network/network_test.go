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
