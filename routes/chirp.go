package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/maximedelporte/chirpy/internal"
)

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, 500, "Something went wrong")
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
	response := cleanBody{cleanResponse}
	internal.RespondWithJSON(w, 200, response)
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
