package types

type contextKeyType string

const UserContextKey contextKeyType = "user"

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}
