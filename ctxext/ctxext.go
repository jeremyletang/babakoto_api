package ctxext

import (
	"context"

	"github.com/jeremyletang/babakoto_api/domain"
)

const (
	accessTokenStringKey = "access_token"
	accessTokenKey       = "domain.AccessToken"
	userKey              = "domain.User"
)

func ExtractAccessTokenString(ctx context.Context) (string, bool) {
	at, ok := ctx.Value(accessTokenStringKey).(string)
	return at, ok
}

func ExtractAccessToken(ctx context.Context) (domain.AccessToken, bool) {
	at, ok := ctx.Value(accessTokenKey).(domain.AccessToken)
	return at, ok
}

func ExtractUser(ctx context.Context) (domain.User, bool) {
	u, ok := ctx.Value(userKey).(domain.User)
	return u, ok
}

func AddAccessTokenString(ctx context.Context, at string) context.Context {
	return context.WithValue(ctx, accessTokenStringKey, at)
}

func AddUser(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func AddAccessToken(ctx context.Context, at domain.AccessToken) context.Context {
	return context.WithValue(ctx, accessTokenKey, at)
}
