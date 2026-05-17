package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
	"github.com/africanecMorj/goods-service_shooeshop/pkg/hash"
	"github.com/africanecMorj/goods-service_shooeshop/pkg/jwtpkg"
)

type AuthService struct {
	Users  *repository.UserRepo
	Tokens *repository.TokenRepo
	Secret []byte
}

func (s *AuthService) Register(ctx context.Context, email, password, role string) error {
	h, err := hash.Hash(password)
	if err != nil {
		return err
	}
	return s.Users.Create(ctx, email, h, role)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.Users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if err = hash.Check(u.Password, password); err != nil {
		return "", "", err
	}

	access, err := jwtpkg.GenerateAccess(s.Secret, u.ID, u.Role)
	if err != nil {
		return "", "", err
	}

	refresh, err := s.newRefresh()
	if err != nil {
		return "", "", err
	}


	err = s.Tokens.Save(ctx, u.ID, refresh, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return "", "", err
	}


	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, oldToken string) (string, string, error) {
	meta, err := s.Tokens.RefreshMeta(ctx, oldToken)
	if err != nil {
		return "", "", err
	}

	access, err := jwtpkg.GenerateAccess(s.Secret, meta.UserID, meta.Role)
	if err != nil {
		return "", "", err
	}

	if !meta.Rotated && time.Now().Before(meta.Exp) {
		return access, oldToken, nil
	}

	newRefresh, err := s.newRefresh()
	if err != nil {
		return "", "", err
	}

	err = s.Tokens.Save(ctx, meta.UserID, newRefresh, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}

func (s *AuthService) newRefresh() (string, error) {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b), nil
}

func (s *AuthService) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	return s.Users.GetUsers(ctx, limit, offset)
}
