package handler

import (
	"net/http"

	"github.com/a-h/templ"
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
		Password: r.FormValue("Password"),
	}
	return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
		InvalidCredentials: "The credentials are invalid",
	}))

	// fmt.Println(credentials)
	// return nil
}
