package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"

)

type ProductRepo struct {
	DB domain.DBTX
}

func (r *ProductRepo) CreateProduct(ctx context.Context, p domain.Product ) (int64, error) {
	var id int64

	err := r.DB.QueryRow(ctx, `
		INSERT INTO products (name, description, price, image_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, p.Name, p.Description, p.Price, p.ImagePath).Scan(&id)

	return id, err
}

func (r *ProductRepo) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	var p domain.Product

	err := r.DB.QueryRow(ctx, `
		SELECT id, name, description, price, image_path, created_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.ImagePath,
		&p.CreatedAt,
	)

	return p, err
}


func (r *ProductRepo) UpdatePartial(ctx context.Context, id int64, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// whitelist to prevent SQL injection
	allowed := map[string]bool{
		"name":        true,
		"description": true,
		"price":       true,
	}

	var sets []string
	var args []interface{}
	i := 1

	for k, v := range fields {
		if !allowed[k] {
			continue
		}
		sets = append(sets, fmt.Sprintf("%s=$%d", k, i))
		args = append(args, v)
		i++
	}

	if len(sets) == 0 {
		return fmt.Errorf("no valid fields provided")
	}

	query := fmt.Sprintf("UPDATE products SET %s WHERE id=$%d",
		strings.Join(sets, ","), i,
	)

	args = append(args, id)

	_, err := r.DB.Exec(ctx, query, args...)
	return err
}

func (r *ProductRepo) GetProducts(
	ctx context.Context,
	limit, offset int,
) ([]domain.Product, int, error) {

	var (
		products []domain.Product
		total    int
	)

	// Total count
	err := r.DB.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM products
	`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Paginated query
	rows, err := r.DB.Query(ctx, `
		SELECT
			id,
			name,
			description,
			price,
			created_at
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Product

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		products = append(products, p)
	}

	// Check iteration errors
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, id int64) (string, error) {
	var imagePath string

	err := r.DB.QueryRow(ctx,
		"DELETE FROM products WHERE id = $1 RETURNING image_path",
		id,
	).Scan(&imagePath)

	if err != nil {
		return "", err
	}

	return imagePath, nil
}

func (r *ProductRepo) PWithTx(tx pgx.Tx) *ProductRepo {
	return &ProductRepo{DB: tx}
}