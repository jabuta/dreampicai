package types

import "github.com/google/uuid"

type contextKeyType string

const UserContextKey contextKeyType = "user"

type AuthenticatedUser struct {
	UserID   uuid.UUID
	Email    string
	LoggedIn bool
}
