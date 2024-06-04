package usecase

import (
	"time"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type TokenUsecase interface {
	New(userID int64, ttl time.Duration, scope string) (*entity.Token, error)
	CheckToken(plaintext string, scope string) (int64, error)
	StartTokenChecker(interval time.Duration)
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
