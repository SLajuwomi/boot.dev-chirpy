package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/slajuwomi/chirpy/internal/auth"
	"github.com/slajuwomi/chirpy/internal/database"
)

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "revoke: failed to get refresh token from header", err)
		return
	}
	dbRefreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "revoke: invalid refresh token", err)
		return
	}
	dbRefreshToken.RevokedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	dbRefreshToken.UpdatedAt = time.Now()
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: dbRefreshToken.RevokedAt,
		Token:     refreshToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "revoke: failed to update refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
