package handlers

import (
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func Removecomment(w http.ResponseWriter, r *http.Request,user database.User) {
    if r.URL.Path == "/removecomment" && r.Method == "GET" {
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
	
		var  content, username string
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
		_, err = db.Exec("DELETE FROM CommentLikes WHERE comment_id = ?", commentID)
		if err != nil {
			http.Error(w, "Failed to delete comment likes/dislikes", http.StatusInternalServerError)
			return
		}
        _, err = db.Exec("DELETE FROM Comments WHERE comment_id = ?", commentID)
        if err != nil {
            http.Error(w, "Failed to delete post", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
