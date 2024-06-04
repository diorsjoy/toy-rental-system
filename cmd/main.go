package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"toy-rental-system/internal/api/handler"
	"toy-rental-system/internal/api/middleware"
	"toy-rental-system/internal/logger"
	"toy-rental-system/internal/mailer"
	"toy-rental-system/internal/repository/postgres"
	"toy-rental-system/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func main() {
	// Setup DB connection
	db, err := sql.Open("postgres", "postgres://postgres:10122004@localhost/toy_rental?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup RabbitMQ connection
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()

	userRepository := postgres.NewUserRepository(db)
	tokenRepository := postgres.NewTokenRepository(db)

	userService := service.NewUserService(userRepository)
	tokenService := service.NewTokenService(tokenRepository)

	mailer := mailer.New("smtp.example.com", 587, "your-email@example.com", "password")

	tokenService.StartTokenChecker(1 * time.Hour)

	r := mux.NewRouter()

	r.HandleFunc("/v1/users", handler.Register(userService, tokenService, mailer)).Methods("POST")
	r.HandleFunc("/v1/users/activated", handler.Activate(userService, tokenService)).Methods("PUT")
	r.HandleFunc("/v1/login", handler.Login(userService, tokenService)).Methods("POST")

	protectedRoutes := r.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware(tokenService))

	logger.Info.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}
