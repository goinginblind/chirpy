package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Fail to encode json: %s\n", err)
		w.WriteHeader(500) //does this actually do anything at this point tho..
		return
	}
}

// handlerValidate handles HTTP requests to validate the length of a "chirp" message.
// It expects a JSON payload with a "body" field and checks if its length does not exceed 140 characters.
// Responds with a JSON indicating validity or an error message if the input is invalid or too long.
func handlerValidate(w http.ResponseWriter, r *http.Request) {
	const maxChirpLen = 140

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Something went wrong"})
		log.Printf("Fail to decode request body: %s\n", err)
		return
	} else if len(params.Body) > maxChirpLen {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Chirp is too long"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"valid": true})
}
