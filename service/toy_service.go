package service

import (
	"toy-rental-system/helpers"
	"toy-rental-system/internal/data"
)

type ToyService interface {
	createToyHandler(toy *data.Toy)
}

type toyService struct {
	toyRepository data.ToyRepository
	helper        helpers.Helpers
}

func NewToyService(repo data.ToyRepository) ToyService {
	return &toyService{
		toyRepository: repo,
	}
}

func (s *toyService) createToyHandler(toy *data.Toy) {

	var inputToy struct {
		Title          string   `json:"title"`
		Description    string   `json:"desc"`
		Details        []string `json:"details,omitempty"`
		Skills         []string `json:"skills"`
		Categories     []string `json:"categories"`
		RecommendedAge string   `json:"recommended_age"`
		Manufacturer   string   `json:"manufacturer"`
		Value          int64    `json:"value"`
		IsAvailable    bool     `json:"is_available"`
		WaitList       []string `json:"wait_list,omitempty"`
	}

	s.helper.



}
