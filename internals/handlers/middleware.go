package handlers

import (
	"fmt"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"strings"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundpath, foundmethod := false, false
		for _, route := range Routes {
			if strings.HasPrefix(route.Path, "/comment") {
				if !IdCheck(w, r) {
					w.WriteHeader(404)
					utils.FileService("error.html", w, Err[404])
					return
				}
			}
			if route.Path == r.URL.Path {
				foundpath = true
				for _, method := range route.Method {
					if r.Method == method {
						foundmethod = true
					}
				}
			}
		}
		if !foundpath {
			w.WriteHeader(404)
			utils.FileService("error.html", w, Err[404])
			return
		}
		if !foundmethod {
			fmt.Println("middleware")
			w.WriteHeader(405)
			utils.FileService("error.html", w, Err[405])
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IdCheck(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/comment") {
		idStr := r.URL.Query().Get("id")
		_, err := strconv.Atoi(idStr)
		if err != nil {
			return false
		}
	}
	return true
}
