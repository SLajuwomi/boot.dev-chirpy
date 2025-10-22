package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	var allChirps []Chirp
	dbAllChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get all chirps from database", err)
		return
	}
	for _, chirp := range dbAllChirps {
		allChirps = append(allChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	sort.Slice(allChirps, func(i, j int) bool {
		return allChirps[i].CreatedAt.Before(allChirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, allChirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	passedUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse passed UUID", err)
		return
	}
	dbChirp, err := cfg.dbQueries.GetSingleChirp(r.Context(), passedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get chirp from datbase", err)
		return
	}
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
