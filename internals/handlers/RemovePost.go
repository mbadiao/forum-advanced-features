package handlers

import (
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func Removepost(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/remove" && r.Method == "GET" {
        id := r.URL.Query().Get("id")
        postID, _ := strconv.Atoi(id)
		if !CheckId(postID) {
			utils.FileService("error",w,Err[400])
			return
		}
        _, err := db.Exec("DELETE FROM Posts WHERE post_id = ?", postID)
        if err != nil {
            http.Error(w, "Failed to delete post", http.StatusInternalServerError)
            return
        }
		_, err = db.Exec("DELETE FROM PostCategories WHERE post_id = ?", postID)
        if err != nil {
            http.Error(w, "Failed to delete post categories", http.StatusInternalServerError)
            return
        }

        _, err = db.Exec("DELETE FROM Comments WHERE post_id = ?", postID)
        if err != nil {
            http.Error(w, "Failed to delete comments", http.StatusInternalServerError)
            return
        }

        _, err = db.Exec("DELETE FROM LikesDislikes WHERE post_id = ?", postID)
        if err != nil {
            http.Error(w, "Failed to delete likes/dislikes", http.StatusInternalServerError)
            return
        }

        _, err = db.Exec("DELETE FROM CommentLikes WHERE comment_id IN (SELECT comment_id FROM Comments WHERE post_id = ?)", postID)
        if err != nil {
            http.Error(w, "Failed to delete comment likes/dislikes", http.StatusInternalServerError)
            return
        }


        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
