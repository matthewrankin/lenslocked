package context

import (
	"context"

	"github.com/matthewrankin/lenslocked/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

// WithUser adds the user to the context.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User gets the user from the given context.
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
