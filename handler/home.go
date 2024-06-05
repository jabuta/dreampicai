package handler

import (
	"net/http"

	"github.com/jabuta/dreampicai/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
