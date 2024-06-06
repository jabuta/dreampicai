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

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Signup().Render(r.Context(), w)
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	signupErrors := auth.SignupErrors{}
	if ok := util.New(params, util.Fields{
		"Email":           util.Rules(util.Min(2), util.Max(50)),
		"Password":        util.Rules(util.Password),
		"ConfirmPassword": util.Rules(util.Equal(params.Password)),
	}).Validate(&signupErrors); !ok {
		fmt.Println("pritning the error")
		return render(r, w, auth.SignupForm(params, signupErrors))
	}

	// if !util.IsValidEmail(params.Email) {
	// 	signupErrors.Email = "Enter a Valid Email"
	// }
	// if reason, ok := util.ValidatePassword(params.Password); !ok {
	// 	signupErrors.Password = reason
	// }
	// if params.Password != params.ConfirmPassword {
	// 	signupErrors.ConfirmPassword = "Passwords do not match"
	// }
	// if len(signupErrors.Email+signupErrors.Password+signupErrors.ConfirmPassword) > 0 {
	// 	fmt.Println("pritning the error")
	// 	return render(r, w, auth.SignupForm(params, signupErrors))
	// }

	sbuser, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		render(r, w, auth.SignupResponse(sbuser.Email, true))
		return err
	}

	return render(r, w, auth.SignupResponse(sbuser.Email, false))
}

func HandleLogInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.LogIn().Render(r.Context(), w)
}

func HandleLogInCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
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

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}
