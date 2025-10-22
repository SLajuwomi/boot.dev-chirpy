package main

import (
	"encoding/json"
	"net/http"

	"github.com/slajuwomi/chirpy/internal/auth"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode parameters ", err)
		return
	}
	curUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get user from database", err)
		return
	}
	valid, err := auth.CheckPasswordHash(params.Password, curUser.HashedPassword)
	if err != nil || !valid {
		respondWithError(w, http.StatusUnauthorized, "Wrong password", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        curUser.ID,
		CreatedAt: curUser.CreatedAt,
		UpdatedAt: curUser.UpdatedAt,
		Email:     curUser.Email,
	})

}
