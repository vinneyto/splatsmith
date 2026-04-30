package services

import (
	"context"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type LoginService struct {
	loginProvider core.LoginProvider
}

func NewLoginService(loginProvider core.LoginProvider) *LoginService {
	return &LoginService{loginProvider: loginProvider}
}

func (s *LoginService) LoginWithPassword(ctx context.Context, username, password string) (core.LoginResult, error) {
	return s.loginProvider.LoginWithPassword(ctx, username, password)
}
