package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jabuta/dreampicai/types"
)

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := types.AuthenticatedUser{}
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		fmt.Println("from the with user middleware")
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
