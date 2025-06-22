package chirps

import (
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/handlers"
	"github.com/google/uuid"
)

func GetChirps(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	// Get query
	query := r.URL.Query()
	authorIDString := query.Get("author_id")
	sort := query.Get("sort")

	// Init variables
	var err error
	var chirps []database.Chirp

	// If author id provided... and if not
	if authorIDString != "" {
		authorUUID, parseErr := uuid.Parse(authorIDString)
		if parseErr != nil {
			log.Printf("Invalid UUID in request: %v\n", parseErr)
			handlers.RespondWithError(w, http.StatusBadRequest, "Incorrect UUID")
			return
		}
		chirps, err = cfg.DB.GetChirpsFiltAuthor(r.Context(), authorUUID)
	} else {
		chirps, err = cfg.DB.GetAllChirps(r.Context())
	}
	if err != nil {
		log.Printf("Failed to get chirps: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Fail to retrieve chirps from the database")
		return
	}

	// Convert to go struct from the db struct, check if there's a 'sort' query
	out := convertManyDBChirps(chirps)
	switch sort {
	case "", "asc":
		// already sorted
	case "desc":
		reverse(out)
	default:
		handlers.RespondWithError(w, http.StatusBadRequest, "Incorrect sorting query parameter")
		return
	}
	handlers.RespondWithJSON(w, http.StatusOK, out)
}

// Small function to make it easier on the eyes on whats happening in the 'case "desc"' above (slice reversed)
func reverse[T any](s []T) {
	for i, j := 0, len(s); i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
