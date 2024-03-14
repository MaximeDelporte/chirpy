package main

import (
	"fmt"
	"github.com/maximedelporte/chirpy/data"
	"github.com/maximedelporte/chirpy/routes"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	cfg := data.ApiConfig{}

	fs := http.FileServer(http.Dir("./static"))
	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", routes.MiddlewareMetricsInc(fsHandler, &cfg))

	mux.HandleFunc("GET /admin/metrics", func(writer http.ResponseWriter, request *http.Request) {
		routes.HandleMetrics(writer, request, &cfg)
	})

	mux.HandleFunc("/api/reset", func(writer http.ResponseWriter, request *http.Request) {
		routes.HandleReset(writer, request, &cfg)
	})

	mux.HandleFunc("GET /api/healthz", routes.HandleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", routes.HandleValidateChirp)

	corsMux := middlewareCors(mux)

	srv := http.Server{
		Addr:              ":8080",
		Handler:           corsMux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Printf("starting server on %s\n", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
