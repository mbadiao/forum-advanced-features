package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func WhoLike(w http.ResponseWriter, r *http.Request, db *sql.DB, post_id int, liked string, disliked string) {
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
	var likestatus bool
	var dislikestatus bool
	err = db.QueryRow("SELECT liked, disliked FROM LikesDislikes WHERE post_id =? and user_id=?", post_id, UserId).Scan(&likestatus, &dislikestatus)
	if err != nil {
		fmt.Println(err.Error())
		utils.FileService("error.html", w, Err[500])
		return
	}

	if liked != "" && likestatus {
		fmt.Println(u.Username, "are", liked, " post", post_id)
	} else if disliked != "" && dislikestatus {
		fmt.Println(u.Username, "are", disliked, " post", post_id)
	}
}
// func NotifTo(w http.ResponseWriter, r *http.Request, user_id int, message string) {
// 	actualcookie := GetCookieHandler(w, r)
// 	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =? and user_id=?", actualcookie, user_id)
// 	if err != nil {
// 		utils.FileService("error.html", w, Err[500])
// 		return
// 	}
// 	w.Write([]byte(message))
// }
