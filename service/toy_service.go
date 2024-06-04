package service

import (
	"fmt"
	"net/http"
	"toy-rental-system/helpers"
	"toy-rental-system/internal/data"
	"toy-rental-system/internal/validator"
)

type ToyService interface {
	createToyHandler(w http.ResponseWriter, r *http.Request)
	showToyHandler(w http.ResponseWriter, r *http.Request)
	listToysHandler(w http.ResponseWriter, r *http.Request)
	updateToyHandler(w http.ResponseWriter, r *http.Request)
	deleteToyHandler(w http.ResponseWriter, r *http.Request)
}

type toyService struct {
	toyRepository data.ToyRepository
}

func (s *toyService) listToysHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title          string
		Skills         []string
		Categories     []string
		RecommendedAge string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Title = helpers.ReadString(qs, "title", "")
	input.Skills = helpers.ReadCSV(qs, "skills", []string{})
	input.Categories = helpers.ReadCSV(qs, "categories", []string{})
	input.Categories = helpers.ReadCSV(qs, "recAge", []string{})
	input.Page = helpers.ReadInt(qs, "page", 1, v)
	input.PageSize = helpers.ReadInt(qs, "page_size", 24, v)
	input.Sort = helpers.ReadString(qs, "sort", "id")
	input.SortSafeList = []string{"title", "skills", "categories", "-title", "-skills", "-categories"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		return
	}

	toys, metadata, err := s.toyRepository.GetAll(input.Title, input.Skills, input.Categories, input.RecommendedAge, input.Filters)
	if err != nil {
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, envelope{"toys": toys, "metadata": metadata}, nil)

}

func (s *toyService) updateToyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIdParam(r)
	if err != nil {
		return
	}

	toy, err := s.toyRepository.Get(id)
	if err != nil {
		return
	}

	var input struct {
		Title          *string   `json:"title"`
		Description    *string   `json:"desc"`
		Details        *[]string `json:"details"`
		Skills         *[]string `json:"skills"`
		Categories     *[]string `json:"categories"`
		RecommendedAge *string   `json:"recommendedAge"`
		Manufacturer   *string   `json:"manufacturer"`
		Value          *int64    `json:"value"`
	}

	err = helpers.ReadJSON(w, r, &input)
	if err != nil {
		return
	}

	if input.Title != nil {
		toy.Title = *input.Title
	}
	if input.Description != nil {
		toy.Description = *input.Description
	}
	if input.Details != nil {
		toy.Details = *input.Details
	}
	if input.Skills != nil {
		toy.Skills = *input.Skills
	}
	if input.Categories != nil {
		toy.Categories = *input.Categories
	}
	if input.RecommendedAge != nil {
		toy.RecommendedAge = *input.RecommendedAge
	}
	if input.Manufacturer != nil {
		toy.Manufacturer = *input.Manufacturer
	}
	if input.Value != nil {
		toy.Value = *input.Value
	}

	v := validator.New()
	if data.ValidateToy(v, toy); !v.Valid() {
		return
	}

	err = s.toyRepository.Update(toy)
	if err != nil {
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, envelope{"toy": toy}, nil)
	if err != nil {
		return
	}

}

func (s *toyService) deleteToyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIdParam(r)
	if err != nil {
		return
	}

	err = s.toyRepository.Delete(id)
	if err != nil {
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, envelope{"message": "Toy deleted successfully"}, nil)

}

func (s *toyService) showToyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIdParam(r)

	if err != nil {
		return
	}

	toy, err := s.toyRepository.Get(id)
	if err != nil {
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, envelope{"toy": toy}, nil)
	if err != nil {
		return
	}
}

func NewToyService(repo data.ToyRepository) ToyService {
	return &toyService{
		toyRepository: repo,
	}
}

type envelope map[string]any

func (s *toyService) createToyHandler(w http.ResponseWriter, r *http.Request) {

	var inputToy struct {
		Title          string   `json:"title"`
		Description    string   `json:"desc"`
		Details        []string `json:"details,omitempty"`
		Skills         []string `json:"skills"`
		Images         []string `json:"images"`
		Categories     []string `json:"categories"`
		RecommendedAge string   `json:"recommended_age"`
		Manufacturer   string   `json:"manufacturer"`
		Value          int64    `json:"value"`
		IsAvailable    bool     `json:"is_available"`
		WaitList       []string `json:"wait_list,omitempty"`
	}

	err := helpers.ReadJSON(w, r, &inputToy)

	if err != nil {
		return
	}

	toy := &data.Toy{
		Title:          inputToy.Title,
		Description:    inputToy.Description,
		Details:        inputToy.Details,
		Skills:         inputToy.Skills,
		Images:         inputToy.Images,
		Categories:     inputToy.Categories,
		RecommendedAge: inputToy.RecommendedAge,
		Manufacturer:   inputToy.Manufacturer,
		Value:          inputToy.Value,
		IsAvailable:    inputToy.IsAvailable,
		WaitList:       inputToy.WaitList,
	}

	v := validator.New()

	if data.ValidateToy(v, toy); !v.Valid() {
		return
	}

	err = s.toyRepository.Insert(toy)
	if err != nil {
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/toys/%d", toy.ID))

	err = helpers.WriteJSON(w, http.StatusCreated, interface{}(envelope{"toy": toy}), headers)
	if err != nil {
		return
	}

}
