package handler

import (
	"net/http"

	"github.com/jabuta/dreampicai/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticateUser(r)
	return render(r, w, settings.Index(user))
}
