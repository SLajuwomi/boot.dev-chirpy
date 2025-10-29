package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/slajuwomi/chirpy/internal/auth"
	"github.com/slajuwomi/chirpy/internal/database"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		// ExpiresInSeconds int    `json:"expires_in_seconds"`
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
	expiresIn := 3600
	// if expiresIn == 0 || expiresIn > 3600 {
	// 	expiresIn = 3600
	// }
	token, err := auth.MakeJWT(curUser.ID, cfg.secret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create jwt", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token", err)
		return
	}
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    curUser.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to add refresh token to database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:           curUser.ID,
		CreatedAt:    curUser.CreatedAt,
		UpdatedAt:    curUser.UpdatedAt,
		Email:        curUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})

}
