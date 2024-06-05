package unit

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"path/filepath"
	"testing"
	_ "toy-rental-system/helpers"
	"toy-rental-system/internal/config"
	"toy-rental-system/internal/data"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository/postgres"
	"toy-rental-system/internal/service"
	"toy-rental-system/internal/validator"
)

func TestDeleteToy(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("DELETE FROM toys WHERE id = ?").WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulating one row affected

	m := &data.ToyModel{DB: db}
	err := m.Delete(1)
	if err != nil {
		t.Errorf("Unexpected error during deletion: %s", err)
	}
}

func TestValidateToyDescription(t *testing.T) {
	v := validator.New()
	toy := data.Toy{
		Title:          "Test",
		Description:    "",
		Skills:         []string{"dasdas", "dasdasd"},
		Images:         []string{"http://", "http://"},
		Categories:     []string{"asdasd", "asdasd"},
		RecommendedAge: "4",
		Manufacturer:   "UK",
		Value:          25,
		IsAvailable:    true,
	}
	data.ValidateToy(v, &toy)
	if v.Valid() {
		t.Errorf("Expected invalid due to description being null")
	}
}

func TestValidateToyImageUrls(t *testing.T) {
	v := validator.New()
	toy := data.Toy{
		Title: "Test",
		Description: `adwdwefwefwefwefdhfsdfkasldfasldfkasdlfaksdhflaskdfalsdfhalsdkfhsdlfksderqweroqweuryqoweirud
dfsdfsdfasdfasdfsdfsdfsdfsdfasdfasdfassdfasdfasdfsdfsadfasdfsdfsdfasakfsdf`,
		Skills:         []string{"dasdas", "dasdasd"},
		Images:         []string{"htp://", "tp://"},
		Categories:     []string{"asdasd", "asdasd"},
		RecommendedAge: "4",
		Manufacturer:   "UK",
		Value:          25,
		IsAvailable:    true,
	}
	data.ValidateToy(v, &toy)
	if v.Valid() {
		t.Errorf("Expected invalid due to links for images being incorrect")
	}
}

func TestSaveSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewSubscriptionRepository(db)
	subscription := &entity.Subscription{
		ID:       1,
		UserID:   1,
		Tokens:   10,
		Price:    5,
		Currency: "KZT",
	}

	mock.ExpectExec(`INSERT INTO subscriptions`).
		WithArgs(subscription.ID, subscription.UserID, subscription.Tokens, subscription.Price, subscription.Currency).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(subscription)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

type MockSubscriptionRepository struct {
	mock.Mock
	cfg config.Config
}

func (m *MockSubscriptionRepository) Save(sub *entity.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

func TestProcessPayment(t *testing.T) {
	s, err := filepath.Abs("toy-rental-system/tests")
	if err != nil {
		t.Errorf("Error reading directory")
	}
	sUp := filepath.Dir(s)
	sUp1 := filepath.Dir(sUp)
	sUp2 := filepath.Dir(sUp1)
	sUp3 := filepath.Dir(sUp2)
	env, err := config.LoadConfig(sUp3)
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	stripeKey := env.StripeSecret
	cfg := &config.Config{
		StripeSecret: stripeKey,
	}
	repo := new(MockSubscriptionRepository)
	subscriptionService := service.NewSubscriptionService(*cfg, repo)

	subscription := &entity.Subscription{
		Price:    1000,
		Currency: "usd",
	}

	err = subscriptionService.ProcessPayment(subscription)
	assert.NoError(t, err)
}

func TestCreateSubscription(t *testing.T) {
	s, err := filepath.Abs("toy-rental-system/tests")
	if err != nil {
		t.Errorf("Error reading directory")
	}
	sUp := filepath.Dir(s)
	sUp1 := filepath.Dir(sUp)
	sUp2 := filepath.Dir(sUp1)
	sUp3 := filepath.Dir(sUp2)
	env, err := config.LoadConfig(sUp3)
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	stripeKey := env.StripeSecret

	cfg := &config.Config{
		StripeSecret: stripeKey,
	}
	repo := new(MockSubscriptionRepository)
	subscriptionService := service.NewSubscriptionService(*cfg, repo)

	subscription := &entity.Subscription{
		ID:       1,
		UserID:   1,
		Tokens:   20,
		Price:    15,
		Currency: "RUB",
	}

	repo.On("Save", subscription).Return(nil)

	err = subscriptionService.Subscribe(subscription)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
