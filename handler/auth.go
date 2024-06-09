package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/jabuta/dreampicai/pkg/sb"
	"github.com/jabuta/dreampicai/pkg/validate"
	"github.com/jabuta/dreampicai/view/auth"
	"github.com/nedpals/supabase-go"
)

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Value:    "",
		Name:     "at",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)

	userToken, err := r.Cookie("at")

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return err
	}
	err = sb.Client.Auth.SignOut(r.Context(), userToken.Value)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

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
	if ok := validate.New(params, validate.Fields{
		"Email":           validate.Rules(validate.Min(2), validate.Max(50)),
		"Password":        validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(validate.Equal(params.Password)),
	}).Validate(&signupErrors); !ok {
		fmt.Println("pritning the error")
		w.WriteHeader(http.StatusUnauthorized)
		return render(r, w, auth.SignupForm(params, signupErrors))
	}

	// if !validate.IsValidEmail(params.Email) {
	// 	signupErrors.Email = "Enter a Valid Email"
	// }
	// if reason, ok := validate.ValidatePassword(params.Password); !ok {
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
		w.WriteHeader(http.StatusUnauthorized)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials are invalid",
		}))
	}

	setAuthCookie(w, resp.AccessToken, resp.ExpiresIn, time.Now().Add(1*time.Hour))
	fmt.Println(credentials)
	redirectCookie, err := r.Cookie("lrd")
	if err != nil {
		return hxRedirect(w, r, "/")
	}
	redirectCookie.MaxAge = -1
	http.SetCookie(w, redirectCookie)
	return hxRedirect(w, r, redirectCookie.Value)
}

func HandleLogInWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:3000/auth/callback/pkce",
		FlowType:   supabase.PKCE,
	})
	if err != nil {
		return err
	}
	setVerifierCookie(w, resp.CodeVerifier)
	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}
	accessInfo, err := sb.DecodeSBJWT(accessToken)
	if err != nil {
		return err
	}
	expiration, err := accessInfo.GetExpirationTime()
	if err != nil {
		return err
	}
	setAuthCookie(w, accessToken, 3600, expiration.Time)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleAuthCallbackPKCE(w http.ResponseWriter, r *http.Request) error {
	authCode := r.URL.Query().Get("code")
	if len(authCode) == 0 {
		return errors.New("no code to parse")
	}
	codeVerifier, err := r.Cookie("cv")
	if err != nil {
		return err
	}
	accessToken, err := sb.Client.Auth.ExchangeCode(r.Context(), supabase.ExchangeCodeOpts{
		AuthCode:     authCode,
		CodeVerifier: codeVerifier.Value,
	})
	if err != nil {
		return err
	}

	setAuthCookie(w, accessToken.AccessToken, accessToken.ExpiresIn, time.Now().Add(1*time.Hour))
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func setAuthCookie(w http.ResponseWriter, accessToken string, maxAge int, expires time.Time) {

	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   maxAge,
		Expires:  expires,
	}
	http.SetCookie(w, cookie)
}

func setVerifierCookie(w http.ResponseWriter, codeVerifier string) {
	cookie := &http.Cookie{
		Value:    codeVerifier,
		Name:     "cv",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   60,
	}
	http.SetCookie(w, cookie)
}
