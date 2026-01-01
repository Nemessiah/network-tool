package interactive

import (
	"strings"
	"testing"
)

func TestPrompt(t *testing.T) {
	tests := []struct {
		name      string
		inputType string
		prompt    string
		input     string
		want      any
	}{
		{
			name:      "test_BasicString",
			inputType: "string",
			prompt:    "enter string",
			input:     "This is a string\n",
			want:      "This is a string",
		},
		{
			name:      "test_Basicint",
			inputType: "int",
			prompt:    "enter int",
			input:     "42\n",
			want:      42,
		},
		{
			name:      "test_Repromptint",
			inputType: "int",
			prompt:    "enter int",
			input:     "This is a string \n42\n",
			want:      42,
		},
		{
			name:      "test_KeyEntry",
			inputType: "string",
			prompt:    "enter string",
			input:     "{{This is a key}}\nThis is not a key\n",
			want:      "This is not a key",
		},
		{
			name:      "test_Panic",
			inputType: "[]string",
			prompt:    "enter int",
			input:     "42\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputType == "[]string" {
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("expected panic, got none")
					}
				}()
			}
			var got any
			input := strings.NewReader(tt.input)
			switch tt.inputType {
			case "int":
				got = Prompt[int](input, tt.prompt)
			case "string":
				got = Prompt[string](input, tt.prompt)
			case "[]string":
				got = Prompt[[]string](input, tt.prompt)
			}

			if got != tt.want {
				t.Fatalf("expected %s, got %d", tt.want, got)
			}
		})
	}
}
