package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jabuta/dreampicai/pkg/sb"
	"github.com/jabuta/dreampicai/pkg/util"
	"github.com/jabuta/dreampicai/view/auth"
	"github.com/nedpals/supabase-go"
)

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func HandleLogInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.LogIn().Render(r.Context(), w)
}

func HandleLogInCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	loginErrors := auth.LoginErrors{}
	if !util.IsValidEmail(credentials.Email) {
		loginErrors.Email = "Enter a Valid Email"
	}
	// if reason, ok := util.ValidatePassword(credentials.Password); !ok {
	// 	loginErrors.Password = reason
	// }
	if len(loginErrors.Email+loginErrors.Password) > 0 {
		return render(r, w, auth.LoginForm(credentials, loginErrors))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials are invalid",
		}))
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println(credentials)
	return nil
}
