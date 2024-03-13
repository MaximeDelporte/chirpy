package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
	FileserverHits int
}

var tmplt *template.Template

func main() {

	mux := http.NewServeMux()

	apiCfg := apiConfig{}

	fs := http.FileServer(http.Dir("./static"))
	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		tmplt, _ = template.ParseFiles("static/metrics.html")
		err := tmplt.Execute(w, apiCfg)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.FileserverHits = 0
		w.WriteHeader(200)
	})

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits += 1
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
