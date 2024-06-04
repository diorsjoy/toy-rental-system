package postgres


import (
	"database/sql"
	"errors"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user *entity.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, tokens) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Tokens)
	return err
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	row := r.db.QueryRow("SELECT id, username, password, tokens FROM users WHERE username = $1", username)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Tokens)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

