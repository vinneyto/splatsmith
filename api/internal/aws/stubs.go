package aws

import (
	"context"
	"fmt"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type loginProviderStub struct{}

func (s *loginProviderStub) LoginWithPassword(context.Context, string, string) (core.LoginResult, error) {
	return core.LoginResult{}, fmt.Errorf("aws login provider: %w", core.ErrNotImplemented)
}
