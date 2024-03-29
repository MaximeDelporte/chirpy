package routes

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/maximedelporte/chirpy/data"
	"github.com/maximedelporte/chirpy/internal/database"

	"github.com/maximedelporte/chirpy/internal"
)

func HandleGetChirps(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps.")
		return
	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	internal.RespondWithJSON(w, http.StatusOK, chirps)
}

func HandleGetChirp(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps.")
		return
	}
	integerId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "Wrong entry.")
		return
	}

	chirp := database.Chirp{ID: -1}

	for _, v := range dbChirps {
		if v.ID == integerId {
			chirp = v
		}
	}

	if chirp.ID == -1 {
		internal.RespondWithError(w, http.StatusNotFound, "Chirp doesn't exist.")
	} else {
		internal.RespondWithJSON(w, http.StatusOK, chirp)
	}
}

func HandleCreateChirp(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Body) >= 140 {
		internal.RespondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanResponse := removeProfane(params.Body)
	type cleanBody struct {
		CleanBody string `json:"cleaned_body"`
	}

	chirp, err := cfg.DB.CreateChirp(cleanResponse)
	if err != nil {
		internal.RespondWithJSON(w, http.StatusInternalServerError, "Couldn't create chirp")
		return

	}

	internal.RespondWithJSON(w, http.StatusCreated, database.Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func removeProfane(body string) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		lWord := strings.ToLower(word)
		if lWord == "kerfuffle" || lWord == "sharbert" || lWord == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
