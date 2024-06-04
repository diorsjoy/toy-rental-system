package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/domain/usecase"
	"toy-rental-system/internal/logger"
	"toy-rental-system/internal/mailer"
	"toy-rental-system/internal/validator"
)

func Register(userUsecase usecase.UserUsecase, tokenUsecase usecase.TokenUsecase, mailer *mailer.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user entity.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Error.Println("Failed to decode user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := userUsecase.Register(&user); err != nil {
			logger.Error.Println("Failed to register user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := tokenUsecase.New(user.ID, 3*24*time.Hour, usecase.ScopeActivation)
		if err != nil {
			logger.Error.Println("Failed to generate activation token:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"userID":          user.ID,
			"activationToken": token.Plaintext,
		}

		plainBody, err := mailer.RenderTemplate("plainBody", data)
		if err != nil {
			logger.Error.Println("Failed to render email template:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		htmlBody, err := mailer.RenderTemplate("htmlBody", data)
		if err != nil {
			logger.Error.Println("Failed to render email template:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = mailer.Send(user.Email, "Welcome to OYNA!", plainBody, htmlBody)
		if err != nil {
			logger.Error.Println("Failed to send email:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": "Successfully registered. Please check your email to verify your account."}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func Activate(userUsecase usecase.UserUsecase, tokenUsecase usecase.TokenUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Token string `json:"token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error.Println("Failed to decode input:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		v := validator.New()
		usecase.ValidateTokenPlaintext(v, input.Token)
		if !v.Valid() {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		userID, err := tokenUsecase.CheckToken(input.Token, usecase.ScopeActivation)
		if err != nil {
			logger.Error.Println("Invalid or expired token:", err)
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}

		if err := userUsecase.Activate(userID); err != nil {
			logger.Error.Println("Failed to activate user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": "Account successfully activated"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func Login(userUsecase usecase.UserUsecase, tokenUsecase usecase.TokenUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error.Println("Failed to decode input:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := userUsecase.Authenticate(input.Email, input.Password)
		if err != nil {
			logger.Error.Println("Invalid email or password:", err)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := tokenUsecase.New(userID, 24*time.Hour, usecase.ScopeAuthentication)
		if err != nil {
			logger.Error.Println("Failed to generate authentication token:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Successfully logged in",
			"token":   token.Plaintext,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
