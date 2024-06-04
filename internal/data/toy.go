package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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
	Images         []string  `json:"image"`
	Categories     []string  `json:"categories"`
	RecommendedAge string    `json:"recommended_age"`
	Manufacturer   string    `json:"manufacturer"`
	Value          int64     `json:"value"`
	IsAvailable    bool      `json:"isAvailable"`
	WaitList       []string  `json:"waitList,omitempty"`
}

func ValidateToy(v *validator.Validator, toy *Toy) {
	v.Check(toy.Title != "", "title", "title must be provided")
	v.Check(len(toy.Title) <= 500, "title", "title must not be more than 500 bytes long")
	v.Check(len(toy.Description) <= 5000, "desc", "Description must not be more than 5000 bytes long")
	v.Check(len(toy.Details) <= 5, "details", "details must not be more than 5")
	v.Check(v.ImageUrlsCheck(toy.Images), "image", "some of image urls is wrong")
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

type ToyRepository interface {
	Insert(toy *Toy) error
	Get(id int64) (*Toy, error)
	Update(toy *Toy) error
	Delete(id int64) error
	GetAll(title string, skills []string, categories []string, recAge string, filters Filters) ([]*Toy, Metadata, error)
}

func (t ToyModel) Insert(toy *Toy) error {
	query := `
INSERT INTO toys (title, desc, details, skills, categories, images, recommended_age, manufacturer, value, is_available, wait_list)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, created_at`

	args := []any{toy.Title, toy.Description, pq.Array(toy.Details), pq.Array(toy.Skills), pq.Array(toy.Categories), pq.Array(toy.Images), toy.RecommendedAge, toy.Manufacturer, toy.Value, toy.IsAvailable, pq.Array(toy.WaitList)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return t.DB.QueryRowContext(ctx, query, args...).Scan(&toy.ID, &toy.CreatedAt)
}

func (t ToyModel) Get(id int64) (*Toy, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, created_at, title, desc, details ,skills, categories, images, recommended_age, manufacturer, value, is_available, wait_list
FROM toys
WHERE id = $1
`
	var toy Toy

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&toy.ID,
		&toy.CreatedAt,
		&toy.Title,
		&toy.Description,
		pq.Array(&toy.Details),
		pq.Array(&toy.Skills),
		pq.Array(&toy.Categories),
		pq.Array(&toy.Images),
		&toy.RecommendedAge,
		&toy.Manufacturer,
		&toy.Value,
		&toy.IsAvailable,
		pq.Array(&toy.WaitList),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &toy, nil
}

func (t ToyModel) Update(toy *Toy) error {

	query := `UPDATE toys
SET title = $1, desc = $2, details = $3, skills = $4, categories = $5, images = $6, recommended_age = $7, manufacturer = $8, value = $9, is_available = $10, wait_list = $11
WHERE id = $12
RETURNING id
`
	args := []any{
		toy.Title,
		toy.Description,
		pq.Array(toy.Details),
		pq.Array(toy.Skills),
		pq.Array(toy.Images),
		pq.Array(toy.Categories),
		toy.RecommendedAge,
		toy.Manufacturer,
		toy.Value,
		toy.IsAvailable,
		toy.WaitList,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&toy.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}

	}
	return nil

}

func (t ToyModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM toys
WHERE id = 1$
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

func (t ToyModel) GetAll(title string, skills []string, categories []string, recAge string, filters Filters) ([]*Toy, Metadata, error) {
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, desc, details, skills, categories, recommended_age, manufacturer, value, is_available, wait_list
FROM toys
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (skills @> $2 OR $2 = '{}')
AND (categories @> $3 OR $3 = '{}')
AND (recommended_age @> $4 OR $4 = '{}')
ORDER BY %s %s, id ASC
LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(skills), pq.Array(categories), pq.Array(recAge), filters.limit(), filters.offset()}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}

	defer rows.Close()

	totalRecords := 0

	toys := []*Toy{}

	for rows.Next() {
		var toy Toy

		err := rows.Scan(
			&totalRecords,
			&toy.ID,
			&toy.CreatedAt,
			&toy.Title,
			&toy.Description,
			pq.Array(&toy.Details),
			pq.Array(&toy.Skills),
			pq.Array(&toy.Categories),
			&toy.RecommendedAge,
			&toy.Manufacturer,
			&toy.Value,
			&toy.IsAvailable,
			pq.Array(&toy.WaitList),
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		toys = append(toys, &toy)

	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.

	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return toys, metadata, nil

}
