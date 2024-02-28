package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func WhoLike(w http.ResponseWriter, r *http.Request, db *sql.DB, post_id int, action string) {
	actualcookie := GetCookieHandler(w, r)
	var UserId int
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =?", actualcookie).Scan(&UserId)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	u := database.User{}
	err = db.QueryRow("SELECT username FROM Users WHERE user_id =?", UserId).Scan(&u.Username)
	if err != nil {
		utils.FileService("error.html", w, Err[500])
		return
	}
	fmt.Println(u.Username, "are", action, " post", post_id)
}
