package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
)

type UserRepo struct{ 
	DB domain.DBTX
}

func (r *UserRepo) Create(ctx context.Context, email, password, role string) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO users(email,password,role) VALUES($1,$2,$3)", email, password, role)
	return err
}

func (r *UserRepo) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	var (
		users []domain.User
		total int
	)

	// Get total count
	err := r.DB.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM users
	`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated users
	rows, err := r.DB.Query(ctx, `
		SELECT id, email, password, role
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User

		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Password,
			&u.Role,
		); err != nil {
			return nil, 0, err
		}

		users = append(users, u)
	}

	// Important: check iteration error
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}


func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRow(ctx, "SELECT id,email,password,role FROM users WHERE email=$1", email).
		Scan(&u.ID, &u.Email, &u.Password,&u.Role,)
	return &u, err
}

func (r *UserRepo) GetByToken(ctx context.Context, token string) (*domain.User, error) {
	var u domain.User
	err := r.DB.QueryRow(ctx, "SELECT t1.* FROM users t1 JOIN refresh_tokens t2 ON t1.id=t2.user_id WHERE t2.token=$1", token).
		Scan(&u.ID, &u.Email, &u.Password,&u.Role,)
	return &u, err
}

func (r *UserRepo) UpdatePartial(ctx context.Context, id int64, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// whitelist to prevent SQL injection
	allowed := map[string]bool{
		"email":    true,
		"role":     true,
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

	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d",
		strings.Join(sets, ","), i,
	)

	args = append(args, id)

	_, err := r.DB.Exec(ctx, query, args...)
	return err
}

