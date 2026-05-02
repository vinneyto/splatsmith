package core

import "context"

type authClaimsCtxKey struct{}

func WithAuthClaims(ctx context.Context, claims AuthClaims) context.Context {
	return context.WithValue(ctx, authClaimsCtxKey{}, claims)
}

func AuthClaimsFromContext(ctx context.Context) (AuthClaims, bool) {
	claims, ok := ctx.Value(authClaimsCtxKey{}).(AuthClaims)
	if !ok {
		return AuthClaims{}, false
	}
	return claims, true
}
