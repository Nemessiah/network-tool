package internal

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defaultConfig []byte

type Config struct {
	Additionalfiles []string `yaml:"additionalfiles"`
	ConfigVersion   string   `yaml:"configversion"`
}

type Template struct {
	Devicetype []string
	Keys       []string
}

var DefaultTemplate = Template{
	Devicetype: []string{"firewall", "switch", "api", "dhcp"},
	Keys:       []string{"interface", "ipaddress", "vlanid", "vlanname", "zone", "cidr", "mask"},
}

type FeatureCommands struct {
	Feature string              `yaml:"feature"`
	Actions map[string][]string `yaml:"actions"`
}

type Deviceinfo struct {
	Devicetype string            `yaml:"devicetype"`
	Commands   []FeatureCommands `yaml:"commands"` // ORDERED
}

type Fullconfig struct {
	Config  Config                `yaml:"config"`
	Vendors map[string]Deviceinfo `yaml:",inline"`
}

type LoadResult struct {
	Config   Fullconfig
	Warnings []string
}

var Reg = regexp.MustCompile(`\{\{([\w-]+)\}\}`)

func LoadConfig(file string, visited map[string]bool) (LoadResult, error) {
	var result LoadResult
	if visited == nil {
		visited = make(map[string]bool)
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return result, err
	}
	file = filepath.Clean(abs)

	// If we've seen this file already, skip it
	if visited[file] {
		result.Warnings = append(
			result.Warnings,
			fmt.Sprintf("skipping already-loaded config %s", file),
		)
		return result, nil
	}
	visited[file] = true

	// Load file from OS
	raw, err := os.ReadFile(file)
	if err != nil {
		return result, err
	}

	// Unmarshal YAML
	err = yaml.Unmarshal(raw, &result.Config)
	if err != nil {
		return result, err
	}

	// Load additional config files if any
	if len(result.Config.Config.Additionalfiles) > 0 {
		for _, path := range result.Config.Config.Additionalfiles {

			path := strings.TrimSpace(path)
			if path == "" {
				result.Warnings = append(result.Warnings, fmt.Sprintf("empty additionalfiles entry in %s", file))
				continue
			}
			if !filepath.IsAbs(path) {
				path = filepath.Join(filepath.Dir(file), path)
			}
			path = filepath.Clean(path)

			additionalCfg, err := LoadConfig(path, visited)
			if err != nil {
				return result, fmt.Errorf("failed to load additional config %q: %w", path, err)
			}
			result.Warnings = append(result.Warnings, additionalCfg.Warnings...)
			// Merge Vendors from additional config
			if result.Config.Vendors == nil {
				result.Config.Vendors = make(map[string]Deviceinfo)
			}
			for k, v := range additionalCfg.Config.Vendors {
				if _, exists := result.Config.Vendors[k]; exists {
					result.Warnings = append(
						result.Warnings,
						fmt.Sprintf(
							"vendor %q already defined, overriding with values from %s",
							k, path,
						),
					)
				}

				result.Config.Vendors[k] = v
			}
		}
	}

	return result, nil
}

func ConfigCheck() (string, error) {

	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}

	configPath := filepath.Join(home, ".objectstaging", "config.yaml")

	_, err = os.Stat(configPath)

	if os.IsNotExist(err) {
		fmt.Printf("No config found. Creating default config at %s\n", configPath)

		os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return configPath, err
		}

		err = os.WriteFile(configPath, defaultConfig, 0644)
		if err != nil {
			return configPath, err
		}
	}

	return configPath, nil
}

func ConfigValidation(fullconfig *Fullconfig) error {
	placeholderRe := regexp.MustCompile(`\{\{(\w+)\}\}`)
	crudOps := []string{"create", "read", "update", "delete"}

	for vendorName, vendor := range fullconfig.Vendors {

		// Validate devicetype
		if !slices.Contains(DefaultTemplate.Devicetype, vendor.Devicetype) {
			return fmt.Errorf("vendor %q has invalid devicetype %q", vendorName, vendor.Devicetype)
		}

		if vendor.Commands == nil {
			continue
		}

		newFeatures := make([]FeatureCommands, 0, len(vendor.Commands))

		for idx, feature := range vendor.Commands {
			featureName := strings.TrimSpace(feature.Feature)
			if featureName == "" {
				return fmt.Errorf("vendor %q has a commands entry at index %d with empty feature name", vendorName, idx)
			}

			// Drop invalid CRUD operations
			cleanActions := make(map[string][]string)
			for opName, cmds := range feature.Actions {
				if slices.Contains(crudOps, opName) {
					cleanActions[opName] = cmds
				} else {
					log.Printf(
						"warning: vendor %q feature %q has invalid operation %q; ignoring",
						vendorName, featureName, opName,
					)
				}
			}

			if len(cleanActions) == 0 {
				log.Printf(
					"warning: vendor %q feature %q has no valid CRUD operations; ignoring feature",
					vendorName, featureName,
				)
				continue
			}

			// Validate placeholders
			for opName, cmds := range cleanActions {
				for _, cmd := range cmds {
					matches := placeholderRe.FindAllStringSubmatch(cmd, -1)
					for _, match := range matches {
						key := match[1]
						if !slices.Contains(DefaultTemplate.Keys, key) {
							return fmt.Errorf(
								"vendor %q feature %q operation %q contains invalid placeholder key %q in command: %q",
								vendorName, featureName, opName, key, cmd,
							)
						}
					}
				}
			}

			feature.Feature = featureName
			feature.Actions = cleanActions
			newFeatures = append(newFeatures, feature)
		}

		vendor.Commands = newFeatures
		fullconfig.Vendors[vendorName] = vendor
	}

	return nil
}

func UsedDeviceTypes(cfg Fullconfig) ([]string, error) {
	unique := make(map[string]bool)
	for _, vendor := range cfg.Vendors {
		unique[vendor.Devicetype] = true
	}

	var usedTypes []string
	for dt := range unique {
		usedTypes = append(usedTypes, dt)
	}

	return usedTypes, nil
}
