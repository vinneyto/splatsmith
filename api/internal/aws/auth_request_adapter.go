package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type ALBAuthRequestAdapter struct{}

type oidcClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func NewALBAuthRequestAdapter() *ALBAuthRequestAdapter { return &ALBAuthRequestAdapter{} }

func (a *ALBAuthRequestAdapter) Adapt(ctx context.Context, req core.AuthRequest) (context.Context, core.AuthRequest) {
	if strings.TrimSpace(req.AuthorizationHeader) != "" {
		return ctx, req
	}

	claims := parseOIDCClaims(req.OIDCDataHeader)
	userID := strings.TrimSpace(claims.Sub)
	if userID == "" {
		userID = strings.TrimSpace(req.OIDCIdentityHeader)
	}
	email := strings.TrimSpace(claims.Email)

	ctx = core.WithAuthClaims(ctx, core.AuthClaims{UserID: userID, Email: email})
	req.AuthorizationHeader = "Bearer alb"
	return ctx, req
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
