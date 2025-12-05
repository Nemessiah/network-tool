package commands

import (
	"fmt"
	"reflect"
)

type Network struct {
	Name   string `json:"name"`
	Vlanid int    `json:"vlanid"`
	Subnet string `json:"subnet"`
}

type NetworkParams struct {
	Interface   string  `json:"interface"`
	Network     Network `json:"network"`
	NetworkType string  `json:"networktype"`
	Description string  `json:"description"`
	Zone        string  `json:"zone"`
}

func MakeStructMap(input NetworkParams) (map[string]string, error) {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if !(value.Kind() == reflect.Struct) {
		return _, fmt.Errorf("Error: input is not a struct!")
	}
	valueType := value.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := valueType.Field(i)

		switch fieldValue.Kind() {
		case int:
			val, err := strconv.itoa(fieldValue.int())
		case string:
			val := fieldValue.String()
		default:
			panic("unsupported type")
		}
	}

}

// func ReplaceKeys (input NetworkParams, command string) (string, error){
// 	var output string
// 	var err error

// 	output = command
// 	// half psuedo code follows
// 	for key, value := input {
// 		regkey := fmt.sprintf("{{%s}}",key)
// 		&output, err = regexp.ReplaceString(command, regkey)

// 		if !err == nil  {
// 			return _, fmt.Errorf(
// 				"Error replacing %s in %s:%s"
// 				key, command, err
// 			)
// 		}

// 	}
// 	return output, nil
// }
