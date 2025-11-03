package main

import (
	"fmt"
)

func Renderjson(info NetworkParams) (string, error) {
	var err error

	output := fmt.Sprintf(
		`{"name": "%s", "vid": %d, "prefix": "%s"}`,
		info.Name, info.VLANID, info.Subnet,
	)

	return output, err
}
