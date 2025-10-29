package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/slajuwomi/chirpy/internal/auth"
	"github.com/slajuwomi/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	passedChirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlerDeleteChirp: failed to parse chirp id", err)
		return
	}
	chirp, err := cfg.dbQueries.GetSingleChirp(r.Context(), passedChirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "handlerDeleteChirp: failed to get chirp from database", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "handlerDeleteChirp: failed to get bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "handlerDeleteChirp: invalid token", err)
		return
	}
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "handlerDeleteChirp: not the owner of this chirp", err)
		return
	}
	err = cfg.dbQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     passedChirpUUID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlerDeletChipr: failed to delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
