package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/likecomment" {
		w.WriteHeader(400)
		utils.FileService("error.html", w, Err[400])
		return
	}
	if r.Method == "POST" {
		found := false
		usercorrespondance := 0
		actualcookie := GetCookieHandler(w, r)
		if actualcookie == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {

			err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =?", actualcookie).Scan(&usercorrespondance)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				} else {
					fmt.Println("erreur at like comment handler , with the query")
					w.WriteHeader(500)
					utils.FileService("error.html", w, Err[500])
					return
				}

			}
			if usercorrespondance == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				found = true
			}
		}
		if found {
			commentidstr := r.FormValue("commentidlike")

			commentid, err00 := strconv.Atoi(commentidstr)
			if err00 != nil {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}
			if !CheckIdlike(commentid){
				w.WriteHeader(400)
                utils.FileService("error.html", w, Err[400])
                return
			}
			postouzid := r.FormValue("postouzid")
			postouzidnbr, err01 := strconv.Atoi(postouzid)
			fmt.Println(postouzidnbr,"postid")
			if err01 != nil {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				fmt.Println("3")

				return
			}
			if !CheckId(postouzidnbr) {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}
			actionLike := r.FormValue("likecomment")
			actionDislike := r.FormValue("dislikecomment")
			countage := 0
			query := "SELECT COUNT(*) FROM CommentLikes WHERE user_id = ? AND comment_id = ?"
			fmt.Println(usercorrespondance, commentid)
			err1 := db.QueryRow(query, usercorrespondance, commentid).Scan(&countage)
			if err1 != nil {
				if err1 == sql.ErrNoRows {
					fmt.Println("no rows returned")
					countage = 0
				} else {
					fmt.Println("hello")
					fmt.Println(err1.Error())
					w.WriteHeader(500)
					utils.FileService("error.html", w, Err[500])
					return
				}
			}

			if countage == 0 {
				database.Insert(db, "CommentLikes", "(comment_id, user_id)", commentid, usercorrespondance)
			}

			request := "SELECT liked, disliked FROM CommentLikes WHERE user_id = ? AND comment_id = ?"
			row := db.QueryRow(request, usercorrespondance, commentid)
			var liked, disliked bool
			err2 := row.Scan(&liked, &disliked)
			if err2 != nil {
				if err2 == sql.ErrNoRows {
					fmt.Println("no row found")
					return
				} else {
					w.WriteHeader(500)
					utils.FileService("error.html", w, Err[500])
					return
				}
			}
			if (actionLike != "" && actionDislike != "") || (actionLike == "" && actionDislike == "") {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])

				return
			}
			if actionLike != "" {
				if (!liked && !disliked) || (!liked && disliked) {
					query := "UPDATE CommentLikes SET liked = ?, disliked = ? WHERE user_id = ? AND comment_id = ?"
					_, err := db.Exec(query, true, false, usercorrespondance, commentid)
					if err != nil {
						http.Redirect(w, r, "/login", http.StatusSeeOther)
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/comment?id="+postouzid, http.StatusSeeOther)
					return
				}

				if liked && !disliked {
					query := "UPDATE CommentLikes SET liked = ?, disliked = ? WHERE user_id = ? AND comment_id = ?"
					_, err := db.Exec(query, false, false, usercorrespondance, commentid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/comment?id="+postouzid, http.StatusSeeOther)
					return
				}

				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return

			}

			if actionDislike != "" {
				if (!liked && !disliked) || (liked && !disliked) {
					query := "UPDATE CommentLikes SET liked = ?, disliked = ? WHERE user_id = ? AND comment_id = ?"
					_, err := db.Exec(query, false, true, usercorrespondance, commentid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/comment?id="+postouzid, http.StatusSeeOther)
					return
				}
				if !liked && disliked {
					query := "UPDATE CommentLikes SET liked = ?, disliked = ? WHERE user_id = ? AND comment_id = ?"
					_, err := db.Exec(query, false, false, usercorrespondance, commentid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/comment?id="+postouzid, http.StatusSeeOther)
					return
				}
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				return
			}
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
	}
}

func GetStatusComment(db *sql.DB, status string, commentID int, userID int) string {
	var etat bool
	var etatstr string
	query := "SELECT " + status + " FROM CommentLikes WHERE comment_id = ? AND user_id = ?"
	err := db.QueryRow(query, commentID, userID).Scan(&etat)
	if err != nil {
		if err == sql.ErrNoRows {
			etatstr = "debut"
			return etatstr
		} else {
			fmt.Println("error at getstatuscomment function")
			fmt.Println(err.Error())
			return ""
		}
	}
	if etat {
		etatstr = "true"
	} else {
		etatstr = "false"
	}
	return etatstr
}

func GetNbrStatusComment(db *sql.DB, status string, commentID int) int {
	count := 0
	query := "SELECT COUNT(*) FROM CommentLikes WHERE comment_id = ? AND " + status + " = true"
	err := db.QueryRow(query, commentID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			fmt.Println("Erreur lors de l'exécution de la requête dans GetNbrStatuscomment:", err)
			return 0
		}
	}
	return count
}
