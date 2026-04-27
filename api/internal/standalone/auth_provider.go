package standalone

import (
	"context"

	"github.com/vinneyto/splatra/api/internal/core"
)

type DevAuthProvider struct {
	devToken     string
	devUserID    string
	devUserEmail string
}

func NewDevAuthProvider(cfg Config) *DevAuthProvider {
	return &DevAuthProvider{devToken: cfg.DevToken, devUserID: cfg.DevUserID, devUserEmail: cfg.DevUserEmail}
}

func (p *DevAuthProvider) ValidateToken(_ context.Context, token string) (core.AuthClaims, error) {
	if token == "" || token != p.devToken {
		return core.AuthClaims{}, core.ErrInvalidToken
	}
	return core.AuthClaims{UserID: p.devUserID, Email: p.devUserEmail}, nil
}
