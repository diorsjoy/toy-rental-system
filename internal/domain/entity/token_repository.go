package entity

import (
	"time"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenRepository interface {
	Insert(token *Token) error
	Get(hash []byte, scope string) (int64, error)
	DeleteExpiredTokens() error
}
