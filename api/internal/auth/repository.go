package auth

import (
	"context"

	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, params queries.CreateUserParams) (queries.User, error)
	GetUserByEmail(ctx context.Context, email string) (queries.User, error)
}

type authRepository struct {
	store *db.Store
}

func NewAuthRepository(store *db.Store) AuthRepository {
	return &authRepository{
		store: store,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
	return r.store.CreateUser(ctx, params)
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (queries.User, error) {
	return r.store.GetUserByEmail(ctx, email)
}
