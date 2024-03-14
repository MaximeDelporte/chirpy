package main

import (
	"fmt"
	"github.com/maximedelporte/chirpy/data"
	"github.com/maximedelporte/chirpy/internal/database"
	"github.com/maximedelporte/chirpy/routes"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("./internal/database/database.json")
	if err != nil {
		fmt.Println("CRASH")
		log.Fatal(err)
	}

	cfg := data.ApiConfig{
		FileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
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
	mux.HandleFunc("GET /api/chirps", func(writer http.ResponseWriter, request *http.Request) {
		routes.HandleGetChirps(writer, request, &cfg)
	})
	mux.HandleFunc("POST /api/chirps", func(writer http.ResponseWriter, request *http.Request) {
		routes.HandleCreateChirp(writer, request, &cfg)
	})

	corsMux := middlewareCors(mux)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("starting server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
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
