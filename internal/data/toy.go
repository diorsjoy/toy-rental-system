package data

import (
	"database/sql"
	"strings"
	"time"
	"toy-rental-system/internal/validator"
)

type Toy struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Title          string    `json:"title"`
	Description    string    `json:"desc"`
	Details        []string  `json:"details,omitempty"`
	Skills         []string  `json:"skills"`
	Image          string    `json:"image"`
	Categories     []string  `json:"categories"`
	RecommendedAge string    `json:"recommended_age"`
	Manufacturer   string    `json:"manufacturer"`
	Value          int64     `json:"value"`
	IsAvailable    bool      `json:"isAvailable"`
	WaitList       []string  `json:"waitList,omitempty"`
}

func ValidateToy(v *validator.ValidatorToy, toy *Toy) {
	v.Check(toy.Title != "", "title", "title must be provided")
	v.Check(len(toy.Title) <= 500, "title", "title must not be more than 500 bytes long")
	v.Check(len(toy.Description) <= 5000, "desc", "Description must not be more than 5000 bytes long")
	v.Check(len(toy.Details) <= 5, "details", "details must not be more than 5")
	v.Check(strings.HasPrefix(toy.Image, "http://"), "image", "image url is wrong")
	v.Check(toy.Categories != nil, "categories", "categories must be provided")
	v.Check(toy.Skills != nil, "skills", "skills must be provided")
	v.Check(len(toy.Categories) >= 1, "categories", "at least 1 category")
	v.Check(len(toy.Skills) >= 1, "skills", "at least 1 skill")
	v.Check(len(toy.Categories) <= 7, "categories", "no more than 7 categories")
	v.Check(len(toy.Skills) <= 7, "Skills", "no more than 7 skills")
	v.Check(validator.Unique(toy.Categories), "categories", "categories should not contain duplicate values")
	v.Check(validator.Unique(toy.Skills), "skills", "skills should not contain duplicate values")
	v.Check(toy.RecommendedAge != "", "recAge", "age must be provided")
	v.Check(toy.Manufacturer != "", "manufacturer", "manufacturer must be provided")
	v.Check(toy.Value >= 1000, "value", "toy value must be more than 1000 tenge")
	v.Check(toy.Value <= 150000, "value", "limit of toy's value is 150.000 tenge")
}

type ToyModel struct {
	DB *sql.DB
}
