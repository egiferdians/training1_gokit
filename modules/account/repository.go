package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: logger,
	}
}

func (repo *repo) CreateUser(ctx context.Context, user Account) error {
	sql := `
		INSERT INTO account (id, email, password)
		VALUES ($1, $2, $3)`

	if user.Email == "" || user.Password == "" {
		return errors.New("User data error.")
	}

	_, err := repo.db.ExecContext(ctx, sql, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUser(ctx context.Context, id string) (string, error) {
	var email string
	err := repo.db.QueryRow("SELECT email FROM account WHERE id=$1", id).Scan(&email)

	if err != nil {
		return "", errors.New("User not found")
	}

	return email, nil
}
