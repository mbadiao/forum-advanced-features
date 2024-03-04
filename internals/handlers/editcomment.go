package handlers

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"time"
)

func Editcomment(w http.ResponseWriter, r *http.Request, user database.User) {
	
	id := r.URL.Query().Get("id")
	commentID, _ := strconv.Atoi(id)
	if !CheckIdcomm(commentID) {
		utils.FileService("error.html", w, Err[400])
		return
	}
	rows, err := db.Query("SELECT Comments.comment_id, Comments.content, Users.username "+
		"FROM Comments "+
		"INNER JOIN Users ON Comments.user_id = Users.user_id "+
		"WHERE Comments.comment_id = ?", commentID)
	if err != nil {
		w.WriteHeader(500)
		utils.FileService("error.html", w, Err[500])
		return
	}
	defer rows.Close()

	var content, username string
	for rows.Next() {
		var postID int
		err = rows.Scan(&postID, &content, &username)
		if err != nil {
			utils.FileService("error.html", w, Err[400])
			return
		}
	}
	if err = rows.Err(); err != nil {
		utils.FileService("error.html", w, Err[400])
		return
	}
	if username != user.Username {
		utils.FileService("error.html", w, Err[400])
		return
	}

	if r.Method == "GET" {
		donnees := Data{
			Mypost: commentID,
		}
		utils.FileService("editcomment.html", w, donnees)
		return
	}
	if r.Method == "POST" {
		cookie, err := r.Cookie("ForumCookie")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		var userId int
		err = db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value=?", cookie.Value).Scan(&userId)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			fmt.Println(err.Error())
			return
		}
		var username string
		var firstname string
		var lastname string
		errScanUser := db.QueryRow("SELECT userName, firstname, lastname FROM Users WHERE user_id=?", userId).Scan(&username, &firstname, &lastname)
		if errScanUser != nil {
			fmt.Println(errScanUser.Error())
			return
		}

		comment := database.Comment{
			PostID:     commentID,
			UserID:     userId,
			Username:   username,
			Lastname:   lastname,
			Firstname:  firstname,
			Formatdate: utils.FormatTimeAgo(time.Now()),
			Content:    r.FormValue("comment"),
		}
		_, err = db.Exec("UPDATE Comments SET user_id = ?, userName = ?, firstname = ?, lastname = ?, formatDate = ?, content = ? WHERE comment_id = ?", comment.UserID, comment.Username, comment.Firstname, comment.Lastname, comment.Formatdate, comment.Content, commentID)
		if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
