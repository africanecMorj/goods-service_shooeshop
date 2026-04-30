package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/africanecMorj/goods-service_shooeshop/internal/repository"
	"github.com/africanecMorj/goods-service_shooeshop/pkg/hash"
	"github.com/africanecMorj/goods-service_shooeshop/pkg/jwtpkg"
)

type AuthService struct {
	Users  *repository.UserRepo
	Tokens *repository.TokenRepo
	Secret []byte
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	h, _ := hash.Hash(password)
	return s.Users.Create(ctx, email, h)
}

func (s *AuthService) Login(ctx context.Context, email, password string, role string) (string, string, error) {
	u, err := s.Users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if err = hash.Check(u.Password, password); err != nil {
		return "", "", err
	}

	access, err := jwtpkg.GenerateAccess(s.Secret, u.ID, role)
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

func (s *AuthService) Refresh(ctx context.Context, oldToken string, role string) (string, string, error) {
	meta, err := s.Tokens.RefreshMeta(ctx, oldToken)
	if err != nil {
		return "", "", err
	}

	access, err := jwtpkg.GenerateAccess(s.Secret, meta.UserID, role)
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
