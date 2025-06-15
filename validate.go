package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Fail to encode json: %s\n", err)
		w.WriteHeader(500) //does this actually do anything at this point tho..
		return
	}
}

func replaceProfanity(chirp string) string {
	profane := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	sep := strings.Split(chirp, " ")
	for i, word := range sep {
		if ok := profane[strings.ToLower(word)]; ok {
			sep[i] = "****"
		}
	}
	return strings.Join(sep, " ")
}

// handlerValidate handles HTTP requests to validate the length of a "chirp" message.
// It expects a JSON payload with a "body" field and checks if its length does not exceed 140 characters.
// Responds with a JSON of a chirp, profanity is filtered out or an error message if the input is invalid or too long.
func handlerValidate(w http.ResponseWriter, r *http.Request) {
	const maxChirpLen = 140

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Something went wrong"})
		log.Printf("Fail to decode request body: %s\n", err)
		return
	} else if len(params.Body) > maxChirpLen {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Chirp is too long"})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": replaceProfanity(params.Body)})
}
