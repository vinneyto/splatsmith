package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type ALBAuthProvider struct{}

func NewALBAuthProvider() *ALBAuthProvider { return &ALBAuthProvider{} }

type oidcClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func (p *ALBAuthProvider) ValidateHeaders(_ context.Context, headers map[string]string) (core.AuthClaims, error) {
	oidcData := strings.TrimSpace(headers["X-Amzn-Oidc-Data"])
	oidcIdentity := strings.TrimSpace(headers["X-Amzn-Oidc-Identity"])
	if oidcData == "" && oidcIdentity == "" {
		return core.AuthClaims{}, fmt.Errorf("missing alb oidc headers: %w", core.ErrUnauthorized)
	}

	claims := parseOIDCClaims(oidcData)
	userID := strings.TrimSpace(claims.Sub)
	if userID == "" {
		userID = oidcIdentity
	}
	if userID == "" {
		return core.AuthClaims{}, fmt.Errorf("missing alb identity: %w", core.ErrUnauthorized)
	}

	return core.AuthClaims{UserID: userID, Email: strings.TrimSpace(claims.Email)}, nil
}

func parseOIDCClaims(raw string) oidcClaims {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return oidcClaims{}
	}
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return oidcClaims{}
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return oidcClaims{}
	}
	var claims oidcClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return oidcClaims{}
	}
	return claims
}
