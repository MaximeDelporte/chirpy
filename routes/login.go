package routes

import (
	"encoding/json"
	"net/http"

	"github.com/maximedelporte/chirpy/data"
	"github.com/maximedelporte/chirpy/internal"
	"github.com/maximedelporte/chirpy/internal/database"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, cfg *data.ApiConfig) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserFrom(params.Email)
	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, "user does not exist")
		return
	}

	passwordMatch := cfg.DB.DoesPasswordMatch(user.Password, params.Password)

	if passwordMatch {
		internal.RespondWithJSON(w, http.StatusOK, database.User{
			ID:    user.ID,
			Email: user.Email,
		})
	} else {
		internal.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
	}
}
