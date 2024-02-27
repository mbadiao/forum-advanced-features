package main

import (
	"fmt"
	"forum/internals/handlers"
	"net/http"
	"os"
)


func main() {
	if len(os.Args) == 1 {
		fs := http.FileServer(http.Dir("./web/static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		for _, route := range handlers.Routes {
			http.Handle(route.Path, handlers.ErrorMiddleware(route.Handler))
		}
		fmt.Println("Server running at:") 
		fmt.Println("> Localhost:    \033[34mhttp://localhost" + handlers.Port + "\033[0m")
		fmt.Println("> disconnect:   \033[31mpress Ctrl+C\033[0m")
		http.ListenAndServe(handlers.Port, nil)
	}
}
