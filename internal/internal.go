package internal

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed default_config.yaml
var defaultConfig []byte

type Config struct {
	Additinalfiles []string `json:"additionalfiles"`
}

type Template struct {
	Devicetype []string `json:"devicetype"`
	Keys       []string `json:"keys"`
}
type Commands struct {
	Create []string `json:"create"`
	Read   []string `json:"read"`
	Update []string `json:"update"`
	Delete []string `json:"delete"`
}

type Deviceinfo struct {
	Devicetype string              `json:"devicetype"`
	Commands   map[string]Commands `json:"commands"`
}

type Fullconfig struct {
	Config   Config   `json:"config"`
	Template Template `json:"template"`
	Vendors  map[string]Deviceinfo
}

func LoadConfig(file string) (Fullconfig, error) {
	var configjson map[string]json.RawMessage
	var config Fullconfig

	rawconfig, err := os.ReadFile(file)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(rawconfig, &configjson)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configjson["config"], &config.Config)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configjson["template"], &config.Template)

	if err != nil {
		return config, err
	}

	config.Vendors = make(map[string]Deviceinfo)

	for key, value := range configjson {
		if key == "config" || key == "template" {
			continue
		}
		var device Deviceinfo
		err := json.Unmarshal(value, &device)
		if err != nil {
			return config, err
		}
		config.Vendors[key] = device
	}

	return config, nil
}

func ConfigCheck() (string, error) {

	var err error

	home, _ := os.UserHomeDir()
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
