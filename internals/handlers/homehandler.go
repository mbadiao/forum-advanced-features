package handlers

import (
	"fmt"
	"forum/internals/database"
	"net/http"
)



func HomeHandler(w http.ResponseWriter, r *http.Request) {
	db := database.CreateTable()
	CookieHandler(w, r, db)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	db := database.CreateTable()
	ActualCookie := GetCookieHandler(w, r)
	stmt, err := db.Prepare(`DELETE FROM Sessions WHERE cookie_value =?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = stmt.Exec(ActualCookie)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
