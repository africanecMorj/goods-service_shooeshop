package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepo struct{ DB *pgxpool.Pool }

func (r *TokenRepo) Save(ctx context.Context, userID int, token string, exp time.Time) error {
	_, err := r.DB.Exec(ctx,
		"INSERT INTO refresh_tokens(user_id,token,expires_at) VALUES($1,$2,$3)",
		userID, token, exp,
	)
	return err
}

func (r *TokenRepo) Delete(ctx context.Context, token string) error {
	_, err := r.DB.Exec(ctx, "DELETE FROM refresh_tokens WHERE token=$1", token)
	return err
}

func (r *TokenRepo) Get(ctx context.Context, token string) (int, error) {
	var userID int
	err := r.DB.QueryRow(ctx,
		"SELECT user_id FROM refresh_tokens WHERE token=$1 AND expires_at>now()",
		token,
	).Scan(&userID)
	return userID, err
}

func (r *TokenRepo) Consume(ctx context.Context, token string) (int, error) {
	var userID int

	err := r.DB.QueryRow(ctx, `
		DELETE FROM refresh_tokens
		WHERE token = $1 AND expires_at > now()
		RETURNING user_id
	`, token).Scan(&userID)

	return userID, err
}

type RefreshResult struct {
	UserID  int
	Exp     time.Time
	Rotated bool
}

func (r *TokenRepo) RefreshMeta(ctx context.Context, token string) (RefreshResult, error) {
	var res RefreshResult

	err := r.DB.QueryRow(ctx, `
		WITH sel AS (
			SELECT user_id, expires_at
			FROM refresh_tokens
			WHERE token = $1
		),
		del AS (
			DELETE FROM refresh_tokens
			WHERE token = $1 AND expires_at <= now()
			RETURNING user_id
		)
		SELECT
			COALESCE(sel.user_id, del.user_id),
			sel.expires_at,
			(del.user_id IS NOT NULL)
		FROM sel
		FULL OUTER JOIN del ON true
	`, token).Scan(&res.UserID, &res.Exp, &res.Rotated)

	return res, err
}