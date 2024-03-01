package handlers

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommentData struct {
	Comment []database.Comment
	AllData AllData
}

func CommentPost(w http.ResponseWriter, r *http.Request, query string) (*AllData, string) {
	allData, err1 := getAllcomment(w, r, query)
	if err1 != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return nil, "err" // Return nil for the pointer type
	}
	return &allData, "" // Return a pointer to allData
}

func DisplayComment(w http.ResponseWriter, r *http.Request) CommentData {
	var comment CommentData
	var CommentData []database.Comment
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		// return nil
	}
	query := "SELECT * FROM Posts WHERE post_id=" + idStr
	allData, err1 := CommentPost(w, r, query)
	if err1 == "err" {
		w.WriteHeader(400)
		utils.FileService("error.html", w, Err[400])
		// return
	}

	rows, err := db.Query("SELECT * FROM Comments WHERE post_id=? ORDER BY creation_date DESC", id)
	if err != nil {
		fmt.Println(err.Error())
		// return nil
	}

	for rows.Next() {
		var comment database.Comment
		err = rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Username, &comment.Firstname, &comment.Lastname, &comment.Formatdate, &comment.CreationDate)
		comment.Formatdate = utils.FormatTimeAgo(comment.CreationDate)
		if err != nil {
			fmt.Println(err.Error())
			// return nil
		}
		comment.NbrLike = GetNbrStatusComment(db, "liked", comment.CommentID)
		comment.NbrDislike = GetNbrStatusComment(db, "disliked", comment.CommentID)
		CommentData = append(CommentData, comment)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		// return nil
	}
	comment.Comment = CommentData
	comment.AllData = *allData
	return comment
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if len(r.FormValue("comment")) > 200 {
		http.Redirect(w, r, "/comment?id="+idStr, http.StatusSeeOther)
	}
	if err != nil {
		return
	}
	if !CheckId(id) {
		w.WriteHeader(400)
		utils.FileService("error.html", w, Err[400])
		return
	}
	if strings.TrimSpace(r.FormValue("comment")) == "" && r.Method == "POST" {
		http.Redirect(w, r, "/comment?id="+idStr, http.StatusSeeOther)
	} else if strings.TrimSpace(r.FormValue("comment")) != "" && r.Method == "POST" {
		RecordComment(w, r)
	}
	utils.FileService("comment.html", w, DisplayComment(w, r))
}

func RecordComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

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
		PostID:     id,
		UserID:     userId,
		Username:   username,
		Lastname:   lastname,
		Firstname:  firstname,
		Formatdate: utils.FormatTimeAgo(time.Now()),
		Content:    r.FormValue("comment"),
	}
	http.Redirect(w, r, "/comment?id="+idStr, http.StatusSeeOther)
	WhoLike(w, r, db, comment.PostID, "", "", comment.Username)
	database.Insert(db, "Comments", "(post_id, user_id, userName, firstname, lastname, formatDate, content)", comment.PostID, comment.UserID, comment.Username, comment.Firstname, comment.Lastname, comment.Formatdate, comment.Content)
}

func CheckId(id int) bool {
	var userid int
	err := db.QueryRow("SELECT post_id FROM Posts WHERE post_id=?", id).Scan(&userid)
	return err == nil
}

func CheckIdlike(id int) bool {
	var userid int
	err := db.QueryRow("SELECT comment_id FROM Comments WHERE comment_id=?", id).Scan(&userid)
	return err == nil
}
