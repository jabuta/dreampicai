package handler

import (
	"net/http"

	"github.com/jabuta/dreampicai/view/account"
)

func HandleAccountIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticateUser(r)
	return render(r, w, account.Index(user))
}
