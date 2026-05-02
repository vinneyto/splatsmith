package services

import (
	"context"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type AuthService struct {
	authProvider core.AuthProvider
}

func NewAuthService(authProvider core.AuthProvider) *AuthService {
	return &AuthService{authProvider: authProvider}
}

func (s *AuthService) Authenticate(ctx context.Context, headers map[string]string) (core.UserIdentity, error) {
	claims, err := s.authProvider.ValidateHeaders(ctx, headers)
	if err != nil {
		return core.UserIdentity{}, core.ErrInvalidToken
	}

	return core.UserIdentity{UserID: claims.UserID, Email: claims.Email}, nil
}
