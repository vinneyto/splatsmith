package aws

import (
	"context"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type ALBAuthRequestAdapter struct{}

func NewALBAuthRequestAdapter() *ALBAuthRequestAdapter { return &ALBAuthRequestAdapter{} }

func (a *ALBAuthRequestAdapter) Adapt(ctx context.Context, req core.AuthRequest) (context.Context, core.AuthRequest) {
	if strings.TrimSpace(req.AuthorizationHeader) == "" {
		ctx = WithALBIdentity(ctx, req.OIDCIdentityHeader, req.OIDCDataHeader)
		req.AuthorizationHeader = "Bearer alb"
	}
	return ctx, req
}
