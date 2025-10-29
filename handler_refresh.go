package main

import (
	"net/http"
	"time"

	"github.com/slajuwomi/chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type refreshResponse struct {
		NewAccessToken string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get refresh token from header", err)
		return
	}
	dbRefreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}
	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token has been revoked", err)
		return
	}
	newAccessToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create new access token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, refreshResponse{
		NewAccessToken: newAccessToken,
	})
}
