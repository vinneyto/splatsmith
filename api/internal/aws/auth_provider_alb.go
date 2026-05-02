package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type ctxKey string

const (
	ctxUserIDKey ctxKey = "aws.alb.user_id"
	ctxEmailKey  ctxKey = "aws.alb.email"
)

type ALBAuthProvider struct{}

func NewALBAuthProvider() *ALBAuthProvider { return &ALBAuthProvider{} }

func WithALBIdentity(ctx context.Context, userID, email string) context.Context {
	ctx = context.WithValue(ctx, ctxUserIDKey, strings.TrimSpace(userID))
	ctx = context.WithValue(ctx, ctxEmailKey, strings.TrimSpace(email))
	return ctx
}

func (p *ALBAuthProvider) ValidateToken(ctx context.Context, _ string) (core.AuthClaims, error) {
	userID, _ := ctx.Value(ctxUserIDKey).(string)
	email, _ := ctx.Value(ctxEmailKey).(string)
	if userID == "" {
		return core.AuthClaims{}, fmt.Errorf("missing alb identity: %w", core.ErrUnauthorized)
	}
	return core.AuthClaims{UserID: userID, Email: email}, nil
}
