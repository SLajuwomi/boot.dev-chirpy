package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) resetRequests(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	if cfg.platform == "dev" {
		err := cfg.dbQueries.DeleteAllUsers(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete all users", err)
			return
		}
	} else {
		respondWithError(w, http.StatusForbidden, "You are not an admin", errors.New("not an admin"))
		return
	}

}
