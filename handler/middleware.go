package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jabuta/dreampicai/pkg/db"
	"github.com/jabuta/dreampicai/types"
	"github.com/jackc/pgx/v5/pgtype"
)

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticateUser(r)
		if !user.LoggedIn {
			cookie := &http.Cookie{
				Name:     "lrd",
				Value:    r.URL.Path,
				Path:     "/",
				MaxAge:   60,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		fmt.Println("from the WithAuth middleware")
		next.ServeHTTP(w, r)
	})

}

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		//no idea how much of this logic needs to go somewhere else, but it is not here

		accessToken, err := r.Cookie("at")
		if err != nil || r.URL.Path == "/log-out" {
			next.ServeHTTP(w, r)
			return
		}

		authedUser, err := decodeUserAccessToken(accessToken.Value)
		if err != nil {
			http.Redirect(w, r, "/log-out", http.StatusForbidden)
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("somestuff", authedUser.Username == "")
		if authedUser.Username == "" {
			authedUser.Username, err = getUserNameByID(r, authedUser)
			if err != nil && r.URL.Path != "/account" {
				ctx := context.WithValue(r.Context(), types.UserContextKey, authedUser)
				fmt.Println("from the WithUser middleware")
				http.Redirect(w, r.WithContext(ctx), "/account", http.StatusSeeOther)
				return
			}
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, authedUser)
		fmt.Println("from the WithUser middleware")
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func getUserNameByID(r *http.Request, user types.AuthenticatedUser) (string, error) {
	userAccount, err := db.Conf.DB.GetUser(r.Context(), pgtype.UUID{Bytes: user.UserID, Valid: true})
	if err != nil {
		return "", err
	}
	return userAccount.Username, nil
}
