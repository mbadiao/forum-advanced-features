package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

type Errors struct {
	Code    int
	Message string
}

var Errs = map[int]Errors{
	500: {
		http.StatusInternalServerError,
		http.StatusText(500),
	},
}

func FileService(str string, w http.ResponseWriter, data any) {
	tmpl, err := template.ParseFiles("./web/templates/" + str)
	if err != nil {
		if str != "error.html" {
			w.WriteHeader(500)
			FileService("error.html", w, Errs[500])
			return
		} else {
			fmt.Println("error while parsing the indicated templates")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("error while executing the template")
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}