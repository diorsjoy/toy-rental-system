package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"log"
	"time"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/domain/usecase"
	"toy-rental-system/internal/validator"
)

type TokenService struct {
	TokenRepo entity.TokenRepository
}

func NewTokenService(repo entity.TokenRepository) *TokenService {
	return &TokenService{
		TokenRepo: repo,
	}
}

func (s *TokenService) New(userID int64, ttl time.Duration, scope string) (*entity.Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = s.TokenRepo.Insert(token)
	return token, err
}

func (s *TokenService) CheckToken(plaintext string, scope string) (int64, error) {
	v := validator.New()
	usecase.ValidateTokenPlaintext(v, plaintext)
	if !v.Valid() {
		return 0, errors.New("invalid token")
	}

	hash := sha256.Sum256([]byte(plaintext))
	return s.TokenRepo.Get(hash[:], scope)
}

func (s *TokenService) StartTokenChecker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.TokenRepo.DeleteExpiredTokens()
				if err != nil {
					log.Println("Error deleting expired tokens:", err)
				}
			}
		}
	}()
}

func generateToken(userID int64, ttl time.Duration, scope string) (*entity.Token, error) {
	token := &entity.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
