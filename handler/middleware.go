package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jabuta/dreampicai/pkg/db"
	"github.com/jabuta/dreampicai/types"
	"github.com/jackc/pgx/v5"
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

// this is dogshit, neet to do better someday
func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		accessToken, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		//necessary for logging out users that dont have an account setup
		if r.URL.Path == "/log-out" {
			next.ServeHTTP(w, r)
			return
		}

		authedUser, err := decodeUserAccessToken(accessToken.Value)
		if err != nil {
			http.Redirect(w, r, "/log-out", http.StatusForbidden)
			return
		}

		//check if the account is setup
		if authedUser.Username == "" {
			userAccount, err := db.Conf.DB.GetUser(r.Context(), pgtype.UUID{
				Bytes: authedUser.UserID,
				Valid: true,
			})
			if errors.Is(err, pgx.ErrNoRows) && r.URL.Path != "/account" {
				ctx := context.WithValue(r.Context(), types.UserContextKey, authedUser)
				fmt.Println("from the WithUser middleware - no account setup")
				http.Redirect(w, r.WithContext(ctx), "/account", http.StatusSeeOther)
				return
			} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				fmt.Println("from the WithUser middleware - database read error")
				ctx := context.WithValue(r.Context(), types.UserContextKey, authedUser)
				handleError(w, r.WithContext(ctx), http.StatusInternalServerError, err)
				return
			}
			authedUser.Username = userAccount.Username
		}
		fmt.Println("from the WithUser middleware")
		ctx := context.WithValue(r.Context(), types.UserContextKey, authedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
