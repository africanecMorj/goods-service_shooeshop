package service

import (
	"context"
	"mime/multipart"
	"os"
	"io"
	"net/http"

	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service/utils"
)

type StreamerService struct {
	ImageRepo *repository.ImageRepo
	ProductRepo *repository.ProductRepo
	DB domain.TxBeginner
}

type ImageStream struct {
	Reader      io.ReadCloser
	ContentType string
}

func (s *StreamerService) UpdateImage(
	ctx context.Context,
	id int64,
	file multipart.File,
	header *multipart.FileHeader,
) (string, error) {

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer tx.Rollback(ctx)

	txP := s.ProductRepo.PWithTx(tx)
	txI := s.ImageRepo.IWithTx(tx)

	product, err := txP.GetProduct(ctx, id)
	if err != nil {
		return "", err
	}

	newPath, err := utils.SaveFile(file, header)
	if err != nil {
		return "", err
	}

	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(newPath)
		}
	}()

	if err = txI.UpdateImage(ctx, id, newPath); err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	cleanup = false

	if product.ImagePath != "" {
		_ = os.Remove(product.ImagePath)
	}

	return newPath, nil
}

func (s *StreamerService) GetImageReader(ctx context.Context, id int64) (*ImageStream, error) {
	path, err := s.ImageRepo.GetImagePath(ctx, id)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 512)
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		file.Close()
		return nil, err
	}

	contentType := http.DetectContentType(buf[:n])

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, err
	}

	return &ImageStream{
		Reader:      file,
		ContentType: contentType,
	}, nil
}