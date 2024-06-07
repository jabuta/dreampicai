package handler

import (
	"log/slog"
	"net/http"

	"github.com/jabuta/dreampicai/types"
)

func hxRedirect(w http.ResponseWriter, r *http.Request, to string) error {
	if len(r.Header.Get("HX-request")) > 0 {
		w.Header().Set("HX-redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}

func getAuthenticateUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}

func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
		}
	}
}
