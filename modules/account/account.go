package account

import "context"

type Account struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Repository interface {
	CreateUser(ctx context.Context, user Account) error
	GetUser(ctx context.Context, id string) (string, error)
}
