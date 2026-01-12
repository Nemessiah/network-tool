package internal

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		want       Fullconfig
		wantErr    bool
		errString  string
		wantWarn   bool
		warnString string
	}{
		{
			name:       "basic_good",
			configFile: "./testdata/basic_good.yaml",
			want: Fullconfig{
				Config: Config{
					Additionalfiles: []string{},
					ConfigVersion:   "1.0.0",
				},
				Vendors: map[string]Deviceinfo{
					"paloalto": {
						Devicetype: "firewall",
						Commands: []FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
									"update": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
								},
							},
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"set vlan {{vlanid}} name {{vlanname}}",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "PathDoesNotExist",
			configFile: "badpath.yaml",
			wantErr:    true,
			errString:  "no such file or directory",
		},
		{
			name:       "BadYaml",
			configFile: "./testdata/badyaml.yaml",
			wantErr:    true,
			errString:  "yaml: line 1: did not find expected ',' or ']'",
		},
		{
			name:       "basic_goodAdditionalFile",
			configFile: "./testdata/additionalfile.yaml",
			want: Fullconfig{
				Config: Config{
					Additionalfiles: []string{
						"./basic_good.yaml",
					},
					ConfigVersion: "1.0.0",
				},
				Vendors: map[string]Deviceinfo{
					"paloalto": {
						Devicetype: "firewall",
						Commands: []FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
									"update": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
								},
							},
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"set vlan {{vlanid}} name {{vlanname}}",
									},
								},
							},
						},
					},
					"additionalpaloalto": {
						Devicetype: "firewall",
						Commands: []FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "BadAdditionalYaml",
			configFile: "./testdata/badadditionalfile.yaml",
			wantErr:    true,
			errString:  "yaml: line 1: did not find expected ',' or ']'",
		},
		{
			name:       "DuplicateAdditionalFile",
			configFile: "./testdata/duplicateadditonalfile01.yaml",
			want: Fullconfig{
				Config: Config{
					Additionalfiles: []string{
						"./basic_good.yaml",
						"./additionalfile.yaml",
					},
					ConfigVersion: "1.0.0",
				},
				Vendors: map[string]Deviceinfo{
					"paloalto": {
						Devicetype: "firewall",
						Commands: []FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
									"update": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
								},
							},
							{
								Feature: "vlan",
								Actions: map[string][]string{
									"create": {
										"set vlan {{vlanid}} name {{vlanname}}",
									},
								},
							},
						},
					},
					"additionalpaloalto": {
						Devicetype: "firewall",
						Commands: []FeatureCommands{
							{
								Feature: "interface",
								Actions: map[string][]string{
									"create": {
										"set network interface {{interface}} ip {{ipaddress}}",
									},
								},
							},
						},
					},
				},
			},
			wantErr:    false,
			wantWarn:   true,
			warnString: "skipping already-loaded config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var visited map[string]bool
			got, err := LoadConfig(tt.configFile, visited)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errString != "" && !strings.Contains(err.Error(), tt.errString) {
					t.Fatalf("error = %q; want substring %q", err.Error(), tt.errString)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantWarn {
				if len(got.Warnings) == 0 {
					t.Fatalf("Expected Warn, got nil")
				}
				if tt.warnString != "" {
					for _, warnString := range got.Warnings {
						if strings.Contains(warnString, tt.warnString) {
							return
						}
					}
					t.Fatalf("warn = %q; want substring %q", strings.Join(got.Warnings, "\n"), tt.warnString)
				}
			} else {
				if len(got.Warnings) != 0 {
					t.Fatalf("unexpected warnings:\n%s", strings.Join(got.Warnings, "\n"))
				}
			}
			if !reflect.DeepEqual(got.Config, tt.want) {
				t.Fatalf("got %#v; want %#v", got.Config, tt.want)
			}
		})
	}
}

func TestConfigCheck(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		want       any
		wantErr    bool
		errString  string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		want      any
		wantErr   bool
		errString string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestUsedDeviceTypes(t *testing.T) {
	tests := []struct {
		name      string
		want      any
		wantErr   bool
		errString string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
