package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Api() {
	server := http.NewServeMux()

	server.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		bodyStr := string(body)
		msg := fmt.Sprintf(
			"your string: \"%s\" has %d characters.",
			bodyStr,
			len(bodyStr),
		)

		out := []byte(msg)
		_, err = os.ReadFile("badfilename.txt")

		if err != nil {
			// Since "the app handles checks", this can be whatever you want:
			// - always 500
			// - always 400
			// - or pass through the error text
			w.WriteHeader(403)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
	})

	log.Fatal(http.ListenAndServe(":8080", server))
}
