package standalone

import (
	"context"
	"fmt"
	"strings"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type FixedTokenAuthRequestAdapter struct {
	token string
	user  core.AuthClaims
}

func NewFixedTokenAuthRequestAdapter(token, userID, userEmail string) (*FixedTokenAuthRequestAdapter, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("dev token is empty")
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("dev user id is empty")
	}
	return &FixedTokenAuthRequestAdapter{
		token: token,
		user: core.AuthClaims{
			UserID: userID,
			Email:  strings.TrimSpace(userEmail),
		},
	}, nil
}

func (a *FixedTokenAuthRequestAdapter) Adapt(ctx context.Context, req core.AuthRequest) (context.Context, core.AuthRequest) {
	if strings.TrimSpace(req.AuthorizationHeader) == "" {
		req.AuthorizationHeader = "Bearer " + a.token
	}
	ctx = core.WithAuthClaims(ctx, a.user)
	return ctx, req
}
