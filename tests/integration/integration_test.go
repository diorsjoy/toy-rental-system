package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
	"toy-rental-system/internal/api/handler"
	_ "toy-rental-system/internal/api/handler"
	"toy-rental-system/internal/config"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository/postgres"
	"toy-rental-system/internal/service"
	"toy-rental-system/pkg/jsonlog"
)

type IntegrationTestSuite struct {
	suite.Suite
	app    *http.Server
	db     *sql.DB
	logger jsonlog.Logger
}

func (suite *IntegrationTestSuite) SetupSuite() {
	s, err := filepath.Abs("toy-rental-system/tests")
	if err != nil {
		fmt.Errorf("Error reading directory")
	}
	sUp := filepath.Dir(s)
	sUp1 := filepath.Dir(sUp)
	sUp2 := filepath.Dir(sUp1)
	sUp3 := filepath.Dir(sUp2)
	env, err := config.LoadConfig(sUp3)
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	db, err := openDB(env)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	suite.db = db

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	suite.logger = *logger

	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	userRepo := postgres.NewUserRepository(db)

	subscriptionService := service.NewSubscriptionService(cfg, subscriptionRepo)
	userService := service.NewUserService(userRepo)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	muxRouter := mux.NewRouter()
	handler.NewUserHandler(muxRouter, userService)

	app := &http.Server{
		Addr:         ":4000",
		Handler:      muxRouter,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	suite.app = app

	go func() {
		err := app.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("could not listen on :4000", err)
		}
	}()
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.app.Close()
	suite.db.Close()
}

func (suite *IntegrationTestSuite) TestCreateSubscription() {
	subscription := &entity.Subscription{
		ID:                   1,
		UserID:               1,
		StripeSubscriptionID: "sub_123",
		PlanID:               1,
		Tokens:               100,
		Price:                1000,
		Currency:             "usd",
	}

	body, _ := json.Marshal(subscription)
	req, _ := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.app.Handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var response entity.Subscription
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), subscription.ID, response.ID)
	assert.Equal(suite.T(), subscription.UserID, response.UserID)
	assert.Equal(suite.T(), subscription.StripeSubscriptionID, response.StripeSubscriptionID)
	assert.Equal(suite.T(), subscription.PlanID, response.PlanID)
	assert.Equal(suite.T(), subscription.Tokens, response.Tokens)
	assert.Equal(suite.T(), subscription.Price, response.Price)
	assert.Equal(suite.T(), subscription.Currency, response.Currency)
}

func (suite *IntegrationTestSuite) TestRegisterUser() {
	user := &entity.User{
		Username: "testuser",
		Password: "password",
		Tokens:   10,
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.app.Handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusCreated, rr.Code)
}

func (suite *IntegrationTestSuite) TestLoginUser() {
	user := &entity.User{
		Username: "testuser",
		Password: "password",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.app.Handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func (suite *IntegrationTestSuite) TestCreateToy() {
	toy := &entity.Toy{
		Name:  "Toy1",
		Price: 100,
	}

	body, _ := json.Marshal(toy)
	req, _ := http.NewRequest("POST", "/toy", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.app.Handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusCreated, rr.Code)
}

func (suite *IntegrationTestSuite) TestGetToy() {
	toyID := 1
	req, _ := http.NewRequest("GET", fmt.Sprintf("/toy/%d", toyID), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.app.Handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func openDB(cfg config.Config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg)
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil

}

type configuration struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}
