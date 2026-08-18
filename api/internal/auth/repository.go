package auth

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
)

type AuthRepository struct {
	store *db.Store
}

func NewAuthRepository(store *db.Store) *AuthRepository {
	return &AuthRepository{
		store: store,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
	return r.store.CreateUser(ctx, params)
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (queries.User, error) {
	return r.store.GetUserByEmail(ctx, email)
}
