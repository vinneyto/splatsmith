package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type ALBAuthProvider struct{}

func NewALBAuthProvider() *ALBAuthProvider { return &ALBAuthProvider{} }

func (p *ALBAuthProvider) ValidateToken(ctx context.Context, token string) (core.AuthClaims, error) {
	if !strings.EqualFold(strings.TrimSpace(token), "alb") {
		return core.AuthClaims{}, fmt.Errorf("unexpected alb auth token: %w", core.ErrInvalidToken)
	}
	claims, ok := core.AuthClaimsFromContext(ctx)
	if !ok || strings.TrimSpace(claims.UserID) == "" {
		return core.AuthClaims{}, fmt.Errorf("missing alb identity: %w", core.ErrUnauthorized)
	}
	return core.AuthClaims{UserID: strings.TrimSpace(claims.UserID), Email: strings.TrimSpace(claims.Email)}, nil
}
