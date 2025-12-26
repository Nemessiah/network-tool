package interactive

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var keyValueReg = regexp.MustCompile(`\{\{.*?\}\}`)

func Prompt[T any](prompt string) T {
	reader := bufio.NewReader(os.Stdin)
	var zero T // default zero value of type T

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if keyValueReg.MatchString(input) {
			fmt.Println("Entry can not be match '{{*}}'")
			continue
		}

		switch any(zero).(type) {
		case int:
			val, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid entry, must be an integer.")
				continue
			}
			return any(val).(T)
		case string:
			return any(input).(T)
		default:
			panic("unsupported type")
		}
	}
}
