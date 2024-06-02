package cmd

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"toy-rental-system/internal/api/handler"
	"toy-rental-system/internal/repository"
	"toy-rental-system/internal/service"
)

func main() {
	cfg := config.LoadConfig("config/config.yaml")
	r := mux.NewRouter()

	db, err := sql.Open("postgres", "postgres://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()

	userRepo := repository.NewUserRepository()
	toyRepo := repository.NewToyRepository()
	subscriptionRepo := repository.NewSubscriptionRepository()

	userService := service.NewUserService(userRepo)
	toyService := service.NewToyService(toyRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, userService)

	handler.NewUserHandler(r, userService)
	handler.NewToyHandler(r, toyService)
	handler.NewSubscriptionHandler(r, subscriptionService)

	log.Fatal(http.ListenAndServe(":8080", r))
}
