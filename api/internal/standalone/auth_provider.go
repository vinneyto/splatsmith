package standalone

import (
	"context"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type DevAuthProvider struct {
	devToken     string
	devUserID    string
	devUserEmail string
	devUsername  string
	devPassword  string
}

func NewDevAuthProvider(cfg Config) *DevAuthProvider {
	return &DevAuthProvider{
		devToken:     cfg.DevToken,
		devUserID:    cfg.DevUserID,
		devUserEmail: cfg.DevUserEmail,
		devUsername:  cfg.DevUsername,
		devPassword:  cfg.DevPassword,
	}
}

func (p *DevAuthProvider) ValidateHeaders(_ context.Context, headers map[string]string) (core.AuthClaims, error) {
	authHeader := strings.TrimSpace(headers["Authorization"])
	token := ""
	if authHeader == "" {
		token = p.devToken
	} else {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(strings.TrimSpace(parts[0]), "Bearer") {
			return core.AuthClaims{}, core.ErrInvalidToken
		}
		token = strings.TrimSpace(parts[1])
	}
	if token == "" || token != p.devToken {
		return core.AuthClaims{}, core.ErrInvalidToken
	}
	return core.AuthClaims{UserID: p.devUserID, Email: p.devUserEmail}, nil
}

func (p *DevAuthProvider) LoginWithPassword(_ context.Context, username, password string) (core.LoginResult, error) {
	if username == "" || password == "" {
		return core.LoginResult{}, core.ErrInvalidCredentials
	}
	if username != p.devUsername || password != p.devPassword {
		return core.LoginResult{}, core.ErrInvalidCredentials
	}
	return core.LoginResult{
		Token: p.devToken,
		User: core.UserIdentity{
			UserID: p.devUserID,
			Email:  p.devUserEmail,
		},
	}, nil
}
