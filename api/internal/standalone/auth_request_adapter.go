package standalone

import (
	"context"
	"fmt"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type FixedTokenAuthRequestAdapter struct {
	token string
}

func NewFixedTokenAuthRequestAdapter(token string) (*FixedTokenAuthRequestAdapter, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("dev token is empty")
	}
	return &FixedTokenAuthRequestAdapter{token: token}, nil
}

func (a *FixedTokenAuthRequestAdapter) Adapt(ctx context.Context, req core.AuthRequest) (context.Context, core.AuthRequest) {
	if strings.TrimSpace(req.AuthorizationHeader) == "" {
		req.AuthorizationHeader = "Bearer " + a.token
	}
	return ctx, req
}
