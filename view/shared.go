package view

import (
	"context"
	"go/types"
)

func authenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
