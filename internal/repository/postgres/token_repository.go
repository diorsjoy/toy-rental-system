package postgres

import (
	"database/sql"
	"toy-rental-system/internal/repository"
)

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) repository.TokenRepository {
	return &tokenRepository{db: db}
}
