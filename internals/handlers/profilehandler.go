package handlers

import (
	"forum/internals/utils"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	utils.FileService("profile.html", w, nil)
}
