package main

import (
	"encoding/json"
	"net/http"

	"github.com/slajuwomi/chirpy/internal/auth"
	"github.com/slajuwomi/chirpy/internal/database"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type newInfo struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	passedInfo := newInfo{}
	err := decoder.Decode(&passedInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "updateUser: failed to decode passed info", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "updateUser: failed to get bearer token", err)
		return
	}
	curUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "updateUser: invalid token", err)
		return
	}

	hashedPassedPassword, err := auth.HashPassword(passedInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "updateUser: failed to hash password", err)
		return
	}
	user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          passedInfo.Email,
		HashedPassword: hashedPassedPassword,
		ID:             curUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "updateUser: failed to update user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
