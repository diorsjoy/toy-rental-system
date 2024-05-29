package main

import (
	"database/sql"
	"log"
	"net/http"
	"toy-rental-system/internal/api/handler"
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

	// Initialize repositories
	userRepository := postgres.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepository)

	// Initialize handlers
	r := mux.NewRouter()
	handler.NewUserHandler(r, userService)

	// Start server
	log.Fatal(http.ListenAndServe(":4000", r))
}
