package service

import (
	"context"
	"fmt"

	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/service/utils"
	"github.com/africanecMorj/goods-service_shooeshop/pkg/hash"
)

type UpdateService struct {
	UserRepo *repository.UserRepo
	ChacheRepo *repository.ChacheRepo
}

func (s *UpdateService) ResetPassword (ctx context.Context, email, code, password string) error {
	u, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	stored, err := s.ChacheRepo.GetCode(ctx, u.Email)
	if err != nil {
		return err
	}

	if stored != code {
		return fmt.Errorf("Invalid code")
	}

	h, err := hash.Hash(password)
	if err != nil {
		return err
	}
	fields := map[string]interface{}{
		"password":h,
	}

	return s.UserRepo.UpdatePartial(ctx, int64(u.ID), fields)

}

func (s *UpdateService) PatchUser (ctx context.Context, id int64, fields map[string]interface{}) error {
	return s.UserRepo.UpdatePartial(ctx, id, fields)
}


func (s *UpdateService) StoreCode (ctx context.Context, email string) (string, error) {
	code := utils.GenerateCode()
	return code, s.ChacheRepo.StoreCode(ctx, email, code)
}