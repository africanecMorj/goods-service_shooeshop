package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
)

type UserRepo struct{ DB *pgxpool.Pool }

func (r *UserRepo) Create(ctx context.Context, email, password string) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO users(email,password) VALUES($1,$2)", email, password)
	return err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRow(ctx, "SELECT id,email,password FROM users WHERE email=$1", email).
		Scan(&u.ID, &u.Email, &u.Password)
	return &u, err
}


