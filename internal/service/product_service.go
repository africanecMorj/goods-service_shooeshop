package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
)

type ProductService struct {
	Repo *repository.ProductRepo
}

// --- helpers ---

func saveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s_%s",
		time.Now().UnixNano(),
		uuid.NewString(),
		filepath.Base(header.Filename),
	)

	path := filepath.Join("uploads", filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return path, nil
}

func detectContentType(f *os.File) (string, error) {
	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buff)
	_, err := f.Seek(0, 0)
	return contentType, err
}

// --- business logic ---

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

	path, err := saveFile(file, header)
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

func (s *ProductService) GetImageReader(ctx context.Context, id int64) (*os.File, string, error) {
	path, err := s.Repo.GetImagePath(ctx, id)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}

	contentType, err := detectContentType(file)
	if err != nil {
		file.Close()
		return nil, "", err
	}

	return file, contentType, nil
}

func (s *ProductService) Patch(ctx context.Context, id int64, fields map[string]interface{}) error {
	return s.Repo.UpdatePartial(ctx, id, fields)
}

func (s *ProductService) UpdateImage(
	ctx context.Context,
	id int64,
	file multipart.File,
	header *multipart.FileHeader,
) (string, error) {

	product, err := s.Repo.GetProduct(ctx, id)
	if err != nil {
		return "", err
	}

	newPath, err := saveFile(file, header)
	if err != nil {
		return "", err
	}

	if err := s.Repo.UpdateImage(ctx, id, newPath); err != nil {
		return "", err
	}

	if product.ImagePath != "" {
		_ = os.Remove(product.ImagePath) // best effort cleanup
	}

	return newPath, nil
}

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, int, error) {
	return s.Repo.GetProducts(ctx, limit, offset)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.Repo.DeleteProduct(ctx, id)
}