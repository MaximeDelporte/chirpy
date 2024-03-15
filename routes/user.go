package routes

import (
	"encoding/json"
	"github.com/maximedelporte/chirpy/data"
	"github.com/maximedelporte/chirpy/internal"
	"github.com/maximedelporte/chirpy/internal/database"
	"net/http"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		internal.RespondWithJSON(w, http.StatusInternalServerError, "Couldn't create user")
		return

	}

	internal.RespondWithJSON(w, http.StatusCreated, database.User{
		ID:    user.ID,
		Email: user.Email,
	})
}
