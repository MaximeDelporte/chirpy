package main

import (
	"encoding/json"
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

	mux.Handle("POST /api/validate_chirp", handleValidateChirpRequest())

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

func handleValidateChirpRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		type error struct {
			Error string `json:"error"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)

		if err != nil {
			respBody := error{Error: "Something went wrong"}
			data, err := json.Marshal(respBody)

			if err != nil {
				log.Printf("Error decoding parameters: %s", err)
				w.WriteHeader(500)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			w.WriteHeader(500)
			w.Write(data)
			return
		}

		if len(params.Body) <= 140 {
			type valid struct {
				Valid bool `json:"valid"`
			}
			respBody := valid{Valid: true}
			data, err := json.Marshal(respBody)

			if err != nil {
				log.Printf("Error marshalling JSON: %s", err)
				w.WriteHeader(500)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			w.WriteHeader(200)
			w.Write(data)
		} else {
			respBody := error{Error: "Chirp is too long"}
			data, err := json.Marshal(respBody)

			if err != nil {
				log.Printf("Error marshalling JSON: %s", err)
				w.WriteHeader(500)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			w.WriteHeader(400)
			w.Write(data)
		}
	})
}
