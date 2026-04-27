package core

import (
	"context"
	"strings"
)

type AuthService struct {
	authProvider AuthProvider
}

func NewAuthService(authProvider AuthProvider) *AuthService {
	return &AuthService{authProvider: authProvider}
}

func (s *AuthService) Authenticate(ctx context.Context, authorizationHeader string) (UserIdentity, error) {
	token := extractBearerToken(authorizationHeader)
	if token == "" {
		return UserIdentity{}, ErrUnauthorized
	}

	claims, err := s.authProvider.ValidateToken(ctx, token)
	if err != nil {
		return UserIdentity{}, ErrInvalidToken
	}

	return UserIdentity{UserID: claims.UserID, Email: claims.Email}, nil
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
