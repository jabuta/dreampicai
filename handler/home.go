package handler

import "net/http"

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) {
	home.Index.Render(r.ContentLength)
}
