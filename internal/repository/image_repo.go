package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepo struct {
	DB *pgxpool.Pool
}

func (r *ImageRepo) GetImagePath(ctx context.Context, id int64) (string, error) {
	var path string

	err := r.DB.QueryRow(ctx,
		`SELECT image_path FROM products WHERE id = $1`, id,
	).Scan(&path)

	return path, err
}

func (r *ImageRepo) UpdateImage(ctx context.Context, id int64, path string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE products SET image_path=$1 WHERE id=$2`,
		path, id,
	)
	return err
}