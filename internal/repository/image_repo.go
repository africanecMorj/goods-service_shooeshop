package repository

import (
	"context"
	
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"

	"github.com/jackc/pgx/v5"
)

type ImageRepo struct {
	DB domain.DBTX
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

func (r *ImageRepo) IWithTx(tx pgx.Tx) *ImageRepo {
	return &ImageRepo{DB: tx}
}