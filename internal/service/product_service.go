package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/africanecMorj/goods-service_shooeshop/internal/service/utils"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
)

type ProductService struct {
	Repo *repository.ProductRepo
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

	path, err := utils.SaveFile(file, header)
	if err != nil {
		return 0, err
	}

	p := domain.Product{
		Name:        name,
		Description: desc,
		Price:       price,
		ImagePath:   path,
	}

	return s.Repo.CreateProduct(ctx, p)
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
	return s.Repo.DeleteProduct(ctx, id)
}