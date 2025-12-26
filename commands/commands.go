package commands

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/nemessiah/network-tool/internal"
)

type Network struct {
	VlanName string `json:"vlanname"`
	Vlanid   int    `json:"vlanid"`
	Subnet   string `json:"subnet"`
}

type NetworkParams struct {
	Interface   string  `json:"interface"`
	Network     Network `json:"network"`
	NetworkType string  `json:"networktype"`
	Description string  `json:"description"`
	Zone        string  `json:"zone"`
	IpAddress   string  `json:"ipaddress"`
}

func reflectParams(input reflect.Value) (map[string]string, error) {
	var valStr string
	output := make(map[string]string)

	if input.Kind() == reflect.Ptr {
		input = input.Elem()
	}
	if !(input.Kind() == reflect.Struct) {
		return nil, fmt.Errorf("input is not a struct")
	}
	valueType := input.Type()
	for i := 0; i < input.NumField(); i++ {
		fieldValue := input.Field(i)
		fieldType := valueType.Field(i)

		// Determine placeholder name from json tag
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			// No json tag, skip
			continue
		}

		// Handle nested structs by recursion
		if fieldValue.Kind() == reflect.Struct {
			nested, err := reflectParams(fieldValue)
			if err != nil {
				return nil, err
			}
			// Merge nested results into out
			for k, v := range nested {
				output[k] = v
			}
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			valStr = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valStr = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Bool:
			if fieldValue.Bool() {
				valStr = "true"
			} else {
				valStr = "false"
			}
		default:
			// For now, ignore unsupported kinds rather than panic
			continue

		}
		output[tag] = valStr
	}
	return output, nil
}

func MakeStructMap(input NetworkParams) (map[string]string, error) {
	value := reflect.ValueOf(input)
	output, err := reflectParams(value)
	if err != nil {
		return nil, fmt.Errorf("unable to reflect input: %s", err)
	}
	return output, nil
}

func ReplaceKeys(keymap map[string]string, command string) (string, error) {
	var output string
	var placeholder string

	output = command

	for key, value := range keymap {
		placeholder = fmt.Sprintf("{{%s}}", key)
		output = strings.ReplaceAll(output, placeholder, value)
	}

	if internal.Reg.MatchString(output) {
		var keys []string
		keyMap := internal.Reg.FindAllStringSubmatch(output, -1)
		for _, value := range keyMap {
			if !slices.Contains(keys, value[1]) {
				keys = append(keys, value[1])
			}
		}
		return "", fmt.Errorf(
			"Command contains unresolved placeholder(s) after replacement.\n"+
				"    Keys:%s \n"+
				"    Command: %s",
			strings.Join(keys, ", "),
			output,
		)
	}

	return output, nil
}

func ExtractVendorCommandsForAction(cfg internal.Fullconfig, action string) map[string][]string {
	output := make(map[string][]string)

	for vendor, device := range cfg.Vendors {
		cmds := make([]string, 0)

		if device.Commands == nil {
			continue
		}
		for _, feature := range device.Commands {
			if feature.Actions == nil {
				continue
			}
			slice := feature.Actions[action]
			if len(slice) > 0 {
				cmds = append(cmds, slice...)
			}
		}
		if len(cmds) > 0 {
			output[vendor] = cmds
		}

	}
	return output
}
