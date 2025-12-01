package internal

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defaultConfig []byte

type Config struct {
	Additionalfiles []string `yaml:"additionalfiles"`
}

type Template struct {
	Devicetype []string `yaml:"devicetype"`
	Keys       []string `yaml:"keys"`
}

type Deviceinfo struct {
	Devicetype string                         `yaml:"devicetype"`
	Commands   map[string]map[string][]string `yaml:"commands"`
}

type Fullconfig struct {
	Config   Config                `yaml:"config"`
	Template Template              `yaml:"template"`
	Vendors  map[string]Deviceinfo `yaml:",inline"`
}

func LoadConfig(file string) (Fullconfig, error) {
	var cfg Fullconfig

	raw, err := os.ReadFile(file)
	if err != nil {
		return cfg, err
	}

	// Unmarshal YAML
	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		return cfg, err
	}

	// Load additional config files if any
	if len(cfg.Config.Additionalfiles) > 0 {
		for _, path := range cfg.Config.Additionalfiles {

			path = filepath.Join(filepath.Dir(file), path)
			additionalCfg, err := LoadConfig(path)
			if err != nil {
				return cfg, fmt.Errorf("failed to load additional config %q: %w", path, err)
			}

			// Merge Vendors from additional config
			if cfg.Vendors == nil {
				cfg.Vendors = make(map[string]Deviceinfo)
			}
			for k, v := range additionalCfg.Vendors {
				cfg.Vendors[k] = v
			}
		}
	}

	return cfg, nil
}

// check if a configuration exists, and create it if not
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

		err = os.WriteFile(configPath, defaultConfig, 0644)
		if err != nil {
			return configPath, err
		}
	}

	return configPath, nil
}

// ConfigValidation validates the Fullconfig and cleans invalid operations.
func ConfigValidation(fullconfig *Fullconfig) (bool, error) {
	placeholderRe := regexp.MustCompile(`\{\{(\w+)\}\}`)
	crudOps := []string{"create", "read", "update", "delete"}

	for vendorName, vendor := range fullconfig.Vendors {

		// Validate devicetype
		if !slices.Contains(fullconfig.Template.Devicetype, vendor.Devicetype) {
			return false, fmt.Errorf(
				"vendor %q has invalid devicetype %q",
				vendorName, vendor.Devicetype,
			)
		}

		// Make a new map to store cleaned commands
		newCommands := make(map[string]map[string][]string)

		for itemName, ops := range vendor.Commands {
			// Drop invalid CRUD operations
			cleanOps := make(map[string][]string)
			for opName, cmds := range ops {
				if slices.Contains(crudOps, opName) {
					cleanOps[opName] = cmds
				} else {
					log.Printf(
						"warning: vendor %q item %q has invalid operation %q; ignoring",
						vendorName, itemName, opName,
					)
				}
			}

			if len(cleanOps) == 0 {
				log.Printf(
					"warning: vendor %q item %q has no valid CRUD operations; ignoring item",
					vendorName, itemName,
				)
				continue
			}

			// Validate placeholders
			for opName, cmds := range cleanOps {
				for _, cmd := range cmds {
					matches := placeholderRe.FindAllStringSubmatch(cmd, -1)
					for _, match := range matches {
						key := match[1]
						if !slices.Contains(fullconfig.Template.Keys, key) {
							return false, fmt.Errorf(
								"vendor %q item %q operation %q contains invalid placeholder key %q in command: %q",
								vendorName, itemName, opName, key, cmd,
							)
						}
					}
				}
			}

			newCommands[itemName] = cleanOps
		}

		// Write back cleaned commands
		vendor.Commands = newCommands
		fullconfig.Vendors[vendorName] = vendor
	}

	return true, nil
}
