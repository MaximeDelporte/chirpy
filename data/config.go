package data

import "github.com/maximedelporte/chirpy/internal/database"

type ApiConfig struct {
	FileserverHits int
	DB             *database.DB
}
