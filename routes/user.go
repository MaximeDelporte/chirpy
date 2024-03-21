package routes

import (
	"encoding/json"
	"net/http"

	"github.com/maximedelporte/chirpy/data"
	"github.com/maximedelporte/chirpy/internal"
	"github.com/maximedelporte/chirpy/internal/database"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		internal.RespondWithJSON(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	internal.RespondWithJSON(w, http.StatusCreated, database.User{
		ID:    user.ID,
		Email: user.Email,
	})
}
