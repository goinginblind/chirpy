package chirps

import "github.com/goinginblind/chirpy/internal/database"

func convertDBChirpToResponse(chirp database.Chirp) chirpParams {
	return chirpParams{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}

func convertManyDBChirps(chirps []database.Chirp) []chirpParams {
	out := make([]chirpParams, len(chirps))
	for i, c := range chirps {
		out[i] = convertDBChirpToResponse(c)
	}
	return out
}
