package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/slajuwomi/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string     `json:"event"`
		Data  dataStruct `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "handlerUpgrade: failed to get api key from header", err)
		return
	}
	if cfg.polkaKey != apiKey {
		respondWithError(w, http.StatusUnauthorized, "handlerUpgrade: incorrect api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlerUpgrade: failed to decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlerUpgrade: failed to parse user id", err)
		return
	}
	_, err = cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "handlerUpgrade: couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "handlerUpgrade: failed to upgrade user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
