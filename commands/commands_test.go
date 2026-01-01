package commands

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nemessiah/network-tool/internal"
)

func TestMakeStructMap(t *testing.T) {

	tests := []struct {
		name string
		data NetworkParams
		want map[string]string
	}{
		{
			name: "test1",
			data: NetworkParams{
				Interface: "ae1.1042",
				Network: Network{
					VlanName: "CLI_IT",
					Vlanid:   1024,
					Subnet:   "10.3.42.0/24",
				},
				NetworkType: "Internal",
				Description: "Spudnik IT",
				Zone:        "z-int-cli-it",
				IpAddress:   "10.3.42.1",
			},
			want: map[string]string{
				"interface":   "ae1.1042",
				"vlanname":    "CLI_IT",
				"vlanid":      "1024",
				"subnet":      "10.3.42.0/24",
				"networktype": "Internal",
				"description": "Spudnik IT",
				"zone":        "z-int-cli-it",
				"ipaddress":   "10.3.42.1",
			},
		},
		{
			name: "test2",
			data: NetworkParams{
				Interface: "ae1.1256",
				Network: Network{
					VlanName: "WAN_Primary",
					Vlanid:   1256,
					Subnet:   "162.255.76.19/29",
				},
				NetworkType: "WAN",
				Description: "Primary ISP",
				Zone:        "z-internet",
				IpAddress:   "162.255.76.18",
			},
			want: map[string]string{
				"interface":   "ae1.1256",
				"vlanname":    "WAN_Primary",
				"vlanid":      "1256",
				"subnet":      "162.255.76.19/29",
				"networktype": "WAN",
				"description": "Primary ISP",
				"zone":        "z-internet",
				"ipaddress":   "162.255.76.18",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeStructMap(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestReflectParams(t *testing.T) {
	type testNetworkParams struct {
		NetworkParams `json:"_"`
		ExtraTag      bool `json:"extratag"`
		Ignored       string
		Bad           []string `json:"bad"`
	}
	type locationInfo struct {
		country string `json:"country"`
		state   string `json:"state"`
	}
	type personalInfo struct {
		name     string       `json:"name"`
		age      int          `json:"age"`
		location locationInfo `json:"location"`
	}

	tests := []struct {
		name      string
		data      any
		want      map[string]string
		wantErr   bool
		errSubstr string
	}{
		{
			name: "test_Basic",
			data: NetworkParams{
				Interface: "ae1.1042",
				Network: Network{
					VlanName: "CLI_IT",
					Vlanid:   1024,
					Subnet:   "10.3.42.0/24",
				},
				NetworkType: "Internal",
				Description: "Spudnik IT",
				Zone:        "z-int-cli-it",
				IpAddress:   "10.3.42.1",
			},
			want: map[string]string{
				"interface":   "ae1.1042",
				"vlanname":    "CLI_IT",
				"vlanid":      "1024",
				"subnet":      "10.3.42.0/24",
				"networktype": "Internal",
				"description": "Spudnik IT",
				"zone":        "z-int-cli-it",
				"ipaddress":   "10.3.42.1",
			},
			wantErr: false,
		},
		{
			name: "test_DifferentStruct",
			data: personalInfo{
				name: "casey",
				age:  35,
				location: locationInfo{
					country: "United States",
					state:   "California",
				},
			},
			want: map[string]string{
				"name":    "casey",
				"age":     "35",
				"country": "United States",
				"state":   "California",
			},
			wantErr: false,
		},
		{
			name: "test_NotStruct",
			data: map[string]string{
				"name":    "casey",
				"age":     "35",
				"country": "United States",
				"state":   "California",
			},
			wantErr:   true,
			errSubstr: "input is not a struct",
		},
		{
			name: "test_ExtraTag",
			data: testNetworkParams{
				NetworkParams: NetworkParams{
					Interface: "ae1.1256",
					Network: Network{
						VlanName: "WAN_Primary",
						Vlanid:   1256,
						Subnet:   "162.255.76.19/29",
					},
					NetworkType: "WAN",
					Description: "Primary ISP",
					Zone:        "z-internet",
					IpAddress:   "162.255.76.18",
				},
				ExtraTag: true,
			},
			want: map[string]string{
				"interface":   "ae1.1256",
				"vlanname":    "WAN_Primary",
				"vlanid":      "1256",
				"subnet":      "162.255.76.19/29",
				"networktype": "WAN",
				"description": "Primary ISP",
				"zone":        "z-internet",
				"ipaddress":   "162.255.76.18",
				"extratag":    "true",
			},
			wantErr: false,
		},
		{
			name: "test_Ignored",
			data: testNetworkParams{
				NetworkParams: NetworkParams{
					Interface: "ae1.1256",
					Network: Network{
						VlanName: "WAN_Primary",
						Vlanid:   1256,
						Subnet:   "162.255.76.19/29",
					},
					NetworkType: "WAN",
					Description: "Primary ISP",
					Zone:        "z-internet",
					IpAddress:   "162.255.76.18",
				},
				ExtraTag: false,
				Ignored:  "This Gets Dropped",
			},
			want: map[string]string{
				"interface":   "ae1.1256",
				"vlanname":    "WAN_Primary",
				"vlanid":      "1256",
				"subnet":      "162.255.76.19/29",
				"networktype": "WAN",
				"description": "Primary ISP",
				"zone":        "z-internet",
				"ipaddress":   "162.255.76.18",
				"extratag":    "false",
			},
			wantErr: false,
		},
		{
			name: "test_BadKind",
			data: testNetworkParams{
				NetworkParams: NetworkParams{
					Interface: "ae1.1256",
					Network: Network{
						VlanName: "WAN_Primary",
						Vlanid:   1256,
						Subnet:   "162.255.76.19/29",
					},
					NetworkType: "WAN",
					Description: "Primary ISP",
					Zone:        "z-internet",
					IpAddress:   "162.255.76.18",
				},
				ExtraTag: false,
				Bad:      []string{"this", "gets", "dropped"},
			},
			want: map[string]string{
				"interface":   "ae1.1256",
				"vlanname":    "WAN_Primary",
				"vlanid":      "1256",
				"subnet":      "162.255.76.19/29",
				"networktype": "WAN",
				"description": "Primary ISP",
				"zone":        "z-internet",
				"ipaddress":   "162.255.76.18",
				"extratag":    "false",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdata := reflect.ValueOf(tt.data)
			got := make(map[string]string, 0)
			got, err := reflectParams(testdata)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Fatalf("error = %q; want substring %q", err.Error(), tt.errSubstr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestReplaceKeys(t *testing.T) {
	tests := []struct {
		name              string
		replacementValues map[string]string
		data              string
		want              string
		wantErr           bool
		errSubstr         string
	}{
		{
			name: "test_Basic",
			replacementValues: map[string]string{
				"interface":   "ae1.1256",
				"vlanname":    "WAN_Primary",
				"vlanid":      "1256",
				"subnet":      "162.255.76.19/29",
				"networktype": "WAN",
				"description": "Primary ISP",
				"zone":        "z-internet",
				"ipaddress":   "162.255.76.18",
				"extratag":    "true",
			},
			data:    "set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
			want:    "set network interface aggregate-ethernet ae1.1256 layer3 ip 162.255.76.18 vlan 1256",
			wantErr: false,
		},
		{
			name: "test_UnresolvedKey",
			replacementValues: map[string]string{
				"interface":   "ae1.1256",
				"vlanname":    "WAN_Primary",
				"vlanid":      "1256",
				"subnet":      "162.255.76.19/29",
				"networktype": "WAN",
				"description": "Primary ISP",
				"zone":        "z-internet",
				"ipaddress":   "162.255.76.18",
				"extratag":    "true",
			},
			data:    "set {{network}} interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
			wantErr: true,
			errSubstr: fmt.Sprintf(
				"Command contains unresolved placeholder(s) after replacement.\n" +
					"    Keys:network \n" +
					"    Command: set {{network}} interface aggregate-ethernet ae1.1256 layer3 ip 162.255.76.18 vlan 1256",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ReplaceKeys(tt.replacementValues, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Fatalf("error = %q; want substring %q", err.Error(), tt.errSubstr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestExtractVendorCommandsForAction(t *testing.T) {
	tests := []struct {
		name      string
		data      internal.Fullconfig
		action    string
		want      map[string][]string
		wantErr   bool
		errSubstr string
	}{
		{
			name: "Test_BasicCreate",
			data: internal.Fullconfig{
				Config: internal.Config{
					Additionalfiles: []string{
						"",
					},
				},
				Vendors: map[string]internal.Deviceinfo{
					"PaloAlto": {
						Devicetype: "firewall",
						Commands: []internal.FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
										"set zone {{zone}} network layer3 {{interface}}",
									},
									"update": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
										"set zone {{zone}} network layer3 {{interface}}",
									},
								},
							},
						},
					},
					"Cisco": {
						Devicetype: "switch",
						Commands: []internal.FeatureCommands{
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"read": {
										"show vlan id {{vlanid}}",
									},
									"update": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"delete": {
										"no vlan {{vlanid}}",
									},
								},
							},
						},
					},
				},
			},

			action: "create",
			want: map[string][]string{
				"PaloAlto": {
					"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
					"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
					"set zone {{zone}} network layer3 {{interface}}",
				},
				"Cisco": {
					"vlan {{vlanid}}",
					"name {{vlanname}}",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_BasicUpdate",
			data: internal.Fullconfig{
				Config: internal.Config{
					Additionalfiles: []string{
						"",
					},
				},
				Vendors: map[string]internal.Deviceinfo{
					"PaloAlto": {
						Devicetype: "firewall",
						Commands: []internal.FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
										"set zone {{zone}} network layer3 {{interface}}",
									},
									"update": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
									},
								},
							},
						},
					},
					"Cisco": {
						Devicetype: "switch",
						Commands: []internal.FeatureCommands{
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"read": {
										"show vlan id {{vlanid}}",
									},
									"update": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"delete": {
										"no vlan {{vlanid}}",
									},
								},
							},
						},
					},
				},
			},

			action: "update",
			want: map[string][]string{
				"PaloAlto": {
					"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
					"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
				},
				"Cisco": {
					"vlan {{vlanid}}",
					"name {{vlanname}}",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_BasicRead",
			data: internal.Fullconfig{
				Config: internal.Config{
					Additionalfiles: []string{
						"",
					},
				},
				Vendors: map[string]internal.Deviceinfo{
					"PaloAlto": {
						Devicetype: "firewall",
						Commands: []internal.FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
										"set zone {{zone}} network layer3 {{interface}}",
									},
									"update": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
									},
								},
							},
						},
					},
					"Cisco": {
						Devicetype: "switch",
						Commands: []internal.FeatureCommands{
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"read": {
										"show vlan id {{vlanid}}",
									},
									"update": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"delete": {
										"no vlan {{vlanid}}",
									},
								},
							},
						},
					},
				},
			},

			action: "read",
			want: map[string][]string{
				"Cisco": {
					"show vlan id {{vlanid}}",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_BasicDelete",
			data: internal.Fullconfig{
				Config: internal.Config{
					Additionalfiles: []string{
						"",
					},
				},
				Vendors: map[string]internal.Deviceinfo{
					"PaloAlto": {
						Devicetype: "firewall",
						Commands: []internal.FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
										"set zone {{zone}} network layer3 {{interface}}",
									},
									"update": {
										"set network interface aggregate-ethernet {{interface}} layer3 ip {{ipaddress}} vlan {{vlanid}}",
										"set network interface aggregate-ethernet {{interface}} description {{vlanname}}",
									},
								},
							},
						},
					},
					"Cisco": {
						Devicetype: "switch",
						Commands: []internal.FeatureCommands{
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"read": {
										"show vlan id {{vlanid}}",
									},
									"update": {
										"vlan {{vlanid}}",
										"name {{vlanname}}",
									},
									"delete": {
										"no vlan {{vlanid}}",
									},
								},
							},
						},
					},
				},
			},
			action: "delete",
			want: map[string][]string{
				"Cisco": {
					"no vlan {{vlanid}}",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := ExtractVendorCommandsForAction(tt.data, tt.action)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}

}
