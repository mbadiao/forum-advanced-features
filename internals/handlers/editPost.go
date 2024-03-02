package handlers

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Editpost(w http.ResponseWriter, r *http.Request, userid database.User) {
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
		err := rows.Scan(&postID, &title1, &content, &username)
		if err != nil {
			utils.FileService("error.html", w, Err[400])
			return
		}
	}
	if err := rows.Err(); err != nil {
		utils.FileService("error.html", w, Err[400])
		return
	}
	if username != userid.Username {
		utils.FileService("error.html", w, Err[400])
		return
	}
	if r.Method == "GET" {
		donnees := Data{
			Mypost: postID,
		}
		utils.FileService("editpost.html", w, donnees)
		return
	}

	var PostS database.Post
	if r.Method == "POST" {
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
		CheckboxValues := r.Form["checkbox"]
		title := strings.TrimSpace(r.FormValue("title"))
		thread := strings.TrimSpace(r.FormValue("thread"))
		if len(CheckboxValues) == 0 || title == "" || thread == "" {
			fmt.Println("value vide")
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
		runeCount := utf8.RuneCountInString(title)
		if runeCount > 38 {
			fmt.Println("tile long")
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
		a := utils.Checkcategory(CheckboxValues)
		if !a {
			fmt.Println("category incorrect")
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
		PhotoURL := uploadHandler(w, r)
		if PhotoURL != "NoPhoto" {
			if PhotoURL == "err400" {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}
			if PhotoURL == "err408" {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[1])
				return
			}
			if PhotoURL == "err500" {
				w.WriteHeader(500)
				utils.FileService("error.html", w, Err[500])
				return
			} else {
				PhotoURL = PhotoURL[5:]
			}
		}
		PostS = database.Post{
			UserID:   userid.UserID,
			Title:    title,
			PhotoURL: PhotoURL,
			Content:  thread,
		}
		_, err = db.Exec("UPDATE Posts SET user_id=?, title=?, PhotoURL=?, content=? WHERE post_id=?", PostS.UserID, PostS.Title, PostS.PhotoURL, PostS.Content, postID)
		if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
		CategoriesId := utils.GetCategory(CheckboxValues)
		_, err = db.Exec("DELETE FROM PostCategories WHERE post_id=?", postID)
		if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}

		for _, v := range CategoriesId {
			_, err = db.Exec("INSERT INTO PostCategories (post_id, category_id) VALUES (?, ?)", postID, v)
			if err != nil {
				w.WriteHeader(500)
				utils.FileService("error.html", w, Err[500])
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
