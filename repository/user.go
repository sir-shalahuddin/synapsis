package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sir-shalahuddin/synapsis/model"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Register(ctx context.Context, email string, password string, username string) (string, error) {
	var userID string

	query := `INSERT INTO users (email, password, username, created_at) VALUES ($1, $2, $3, NOW()) RETURNING ID`

	err := r.db.QueryRowxContext(ctx, query, email, password, username).Scan(&userID)

	if err != nil {
		log.Println("User Repository Error Create User : ", err)
		return "", errors.New("failed to create user in the database")
	}

	return userID, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `select * from users where email = $1`

	err := r.db.QueryRowxContext(ctx, query, email).StructScan(&user)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Println("User Repository Error Get User : ", err)
		return nil, errors.New("failed to get user by email in the database")
	}

	return &user, nil
}
