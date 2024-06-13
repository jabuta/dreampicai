package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jabuta/dreampicai/pkg/sb"
	"github.com/jabuta/dreampicai/types"
)

func decodeUserAccessToken(token string) (types.AuthenticatedUser, error) {
	userClaims, err := sb.GetUserClaims(token)
	if err != nil {
		return types.AuthenticatedUser{}, err
	}
	userID, err := uuid.Parse(fmt.Sprint(userClaims.Subject)) //Supabase returns user_ID in subject
	if err != nil {
		return types.AuthenticatedUser{}, err
	}
	var user = types.AuthenticatedUser{
		LoggedIn: fmt.Sprint(userClaims.Audience) == "[authenticated]",
		UserID:   userID,
		Email:    userClaims.Email,
		Username: userClaims.Username,
	}

	return user, nil
}

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
