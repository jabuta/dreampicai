package handler

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jabuta/dreampicai/pkg/sb"
	"github.com/jabuta/dreampicai/types"
	"github.com/jabuta/dreampicai/view/systemerror"
)

func handleError(w http.ResponseWriter, r *http.Request, httpStatus int, callbackError error) {
	log.Println(callbackError)
	w.WriteHeader(httpStatus)
	err := render(r, w, systemerror.ErrorPage(httpStatus))
	if err != nil {
		http.Error(w, "The error function errored", http.StatusInternalServerError)
	}
}

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

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
