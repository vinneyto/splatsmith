package services

import (
	"context"
	"strings"

	"github.com/vinneyto/ariadne/api/internal/core"
)

type AuthService struct {
	authProvider core.AuthProvider
}

func NewAuthService(authProvider core.AuthProvider) *AuthService {
	return &AuthService{authProvider: authProvider}
}

func (s *AuthService) Authenticate(ctx context.Context, authorizationHeader string) (core.UserIdentity, error) {
	token := extractBearerToken(authorizationHeader)
	if token == "" {
		return core.UserIdentity{}, core.ErrUnauthorized
	}

	claims, err := s.authProvider.ValidateToken(ctx, token)
	if err != nil {
		return core.UserIdentity{}, core.ErrInvalidToken
	}

	return core.UserIdentity{UserID: claims.UserID, Email: claims.Email}, nil
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
