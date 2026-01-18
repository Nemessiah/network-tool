package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nemessiah/network-tool/commands"
	"github.com/nemessiah/network-tool/internal"
)

func ApiEntry() {
	server := http.NewServeMux()

	server.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var input commands.NetworkParams
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // catches typos like "networkType" vs "networktype"
		if err := dec.Decode(&input); err != nil {
			http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if input.Action == "" {
			http.Error(w, "missing required field: action", http.StatusBadRequest)
			return
		}
		log.Printf("action bytes=%v", []byte(input.Action))

		// TEMP debug (remove once fixed)
		log.Printf("decoded input: %+v", input)
		log.Printf("config: %+v", internal.AppConfig)

		extractedCommands := commands.ExtractVendorCommandsForAction(internal.AppConfig, input.Action)
		log.Printf("action=%q extracted vendors=%d", input.Action, len(extractedCommands))

		outputCommands := make(map[string][]string)
		reflectedInput, err := commands.MakeStructMap(input)
		if err != nil {
			http.Error(w, "reflect input: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for vendor, commandArray := range extractedCommands {
			replaced := make([]string, 0, len(commandArray))
			for _, cmd := range commandArray {
				temp, err := commands.ReplaceKeys(reflectedInput, cmd)
				if err != nil {
					http.Error(w, "replace keys: "+err.Error(), http.StatusBadRequest)
					return
				}
				replaced = append(replaced, temp)
			}
			outputCommands[vendor] = replaced
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(outputCommands)
	})

	log.Fatal(http.ListenAndServe(":8080", server))
}
