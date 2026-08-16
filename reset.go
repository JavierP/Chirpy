package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden\n"))
		return
	}
	//cfg.fileserverHits.Swap(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not reset\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users have been reset\n"))
}
