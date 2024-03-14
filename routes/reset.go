package routes

import (
	"github.com/maximedelporte/chirpy/data"
	"net/http"
)

func HandleReset(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	cfg.FileserverHits = 0
	w.WriteHeader(200)
}
