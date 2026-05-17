package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/africanecMorj/goods-service_shooeshop/internal/service/utils"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
)

type ProductService struct {
	Repo *repository.ProductRepo
	DB domain.TxBeginner
}


func (s *ProductService) CreateProduct(
	ctx context.Context,
	name, desc string,
	price float64,
	file multipart.File,
	header *multipart.FileHeader,
) (int64, error) {

	if name == "" {
		return 0, fmt.Errorf("name is required")
	}

	if price <= 0 {
		return 0, fmt.Errorf("invalid price")
	}

	defer file.Close()

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	txRepo := s.Repo.PWithTx(tx)

	path, err := utils.SaveFile(file, header)
	if err != nil {
		return 0, err
	}

	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(path)
		}
	}()

	p := domain.Product{
		Name:        name,
		Description: desc,
		Price:       price,
		ImagePath:   path,
	}

	id, err := txRepo.CreateProduct(ctx, p)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	committed = true
	cleanup = false

	return id, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	return s.Repo.GetProduct(ctx, id)
}


func (s *ProductService) Patch(ctx context.Context, id int64, fields map[string]interface{}) error {
	return s.Repo.UpdatePartial(ctx, id, fields)
}

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, int, error) {
	return s.Repo.GetProducts(ctx, limit, offset)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	repoTx := s.Repo.PWithTx(tx)

	imagePath, err := repoTx.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	if err := os.Remove(imagePath); err != nil {
		return err
	}

	return nil
}