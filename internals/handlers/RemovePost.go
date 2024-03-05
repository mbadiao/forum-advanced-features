package handlers

import (
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func Removepost(w http.ResponseWriter, r *http.Request,user database.User) {
    if r.URL.Path == "/remove" && r.Method == "GET" {
		id := r.URL.Query().Get("id")
		postID, _ := strconv.Atoi(id)
		if !CheckId(postID) {
			utils.FileService("error.html", w, Err[400])
			return
		}
		rows, err := db.Query("SELECT Posts.post_id, Posts.title, Posts.content, Users.username "+
			"FROM Posts "+
			"INNER JOIN Users ON Posts.user_id = Users.user_id "+
			"WHERE Posts.post_id = ?", postID)
		if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
		defer rows.Close()
	
		var title1, content, username string
		for rows.Next() {
			var postID int
			err = rows.Scan(&postID, &title1, &content, &username)
			if err != nil {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}
		}
		if err = rows.Err(); err != nil {
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
		if username != user.Username {
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
        _, err = db.Exec("DELETE FROM Posts WHERE post_id = ?", postID)
        if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
        }
		_, err = db.Exec("DELETE FROM PostCategories WHERE post_id = ?", postID)
        if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
        }
		
		_, err = db.Exec("DELETE FROM CommentLikes WHERE comment_id IN (SELECT comment_id FROM Comments WHERE post_id = ?)", postID)
		if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
        _, err = db.Exec("DELETE FROM Comments WHERE post_id = ?", postID)
        if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
        }

        _, err = db.Exec("DELETE FROM LikesDislikes WHERE post_id = ?", postID)
        if err != nil {
            w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
