package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jabuta/dreampicai/pkg/sb"
	"github.com/jabuta/dreampicai/types"
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
		accessToken, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userClaims, err := sb.GetUserClaims(accessToken.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if fmt.Sprint(userClaims.Audience) != "[authenticated]" {
			accessToken.MaxAge = -1
			accessToken.Value = ""
			http.SetCookie(w, accessToken)
			next.ServeHTTP(w, r)
			return
		}
		userID, err := uuid.Parse(fmt.Sprint(userClaims.Subject))
		if err != nil {
			accessToken.MaxAge = -1
			accessToken.Value = ""
			http.SetCookie(w, accessToken)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, types.AuthenticatedUser{
			UserID:   userID,
			Email:    userClaims.Email,
			LoggedIn: true,
		})
		fmt.Println("from the WithUser middleware")
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
