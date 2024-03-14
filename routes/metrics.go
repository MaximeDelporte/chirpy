package routes

import (
	"github.com/maximedelporte/chirpy/data"
	"html/template"
	"net/http"
)

func HandleMetrics(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	tmplt, _ := template.ParseFiles("static/metrics.html")
	err := tmplt.Execute(w, cfg)
	if err != nil {
		return
	}
}

func MiddlewareMetricsInc(next http.Handler, cfg *data.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits += 1
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
