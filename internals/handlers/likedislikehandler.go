package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/likedislike" {
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
					fmt.Println("erreur at like dislike handler , with the query")
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
			postidstr := r.FormValue("postidouz")
			postid, err := strconv.Atoi(postidstr)
			if !CheckId(postid) {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				fmt.Println("2")

				return
			}
			if err != nil {
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				fmt.Println("3")

				return
			}

			actionLike := r.FormValue("actionlike")
			actionDislike := r.FormValue("actiondislike")
			count := 0
			query := "SELECT COUNT(*) FROM LikesDislikes WHERE user_id = ? AND post_id = ?"
			err1 := db.QueryRow(query, usercorrespondance, postid).Scan(&count)
			if err1 != nil {
				if err1 == sql.ErrNoRows {
					fmt.Println("no rows returned")
					count = 0
				} else {
					w.WriteHeader(500)
					utils.FileService("error.html", w, Err[500])
					return
				}
			}

			if count == 0 {
				database.Insert(db, "LikesDislikes", "(post_id, user_id)", postid, usercorrespondance)
			}

			request := "SELECT liked, disliked FROM LikesDislikes WHERE user_id = ? AND post_id = ?"
			row := db.QueryRow(request, usercorrespondance, postid)
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
				fmt.Println("c'est mort")
				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])

				return
			}
			if actionLike != "" {
				if (!liked && !disliked) || (!liked && disliked) {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, true, false, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

				if liked && !disliked {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, false, false, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

				w.WriteHeader(400)
				utils.FileService("error.html", w, Err[400])
				fmt.Println("4")

				return

			}

			if actionDislike != "" {
				if (!liked && !disliked) || (liked && !disliked) {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, false, true, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
				if !liked && disliked {
					query := "UPDATE LikesDislikes SET liked = ?, disliked = ? WHERE user_id = ? AND post_id = ?"
					_, err := db.Exec(query, false, false, usercorrespondance, postid)
					if err != nil {
						fmt.Println("error to uptade like or dislike")
						return
					}
					http.Redirect(w, r, "/", http.StatusSeeOther)
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

func GetStatus(db *sql.DB, status string, postID int, userID int) string {
	var etat bool
	var etatstr string
	query := "SELECT " + status + " FROM LikesDislikes WHERE post_id = ? AND user_id = ?"
	err := db.QueryRow(query, postID, userID).Scan(&etat)
	if err != nil {
		if err == sql.ErrNoRows {
			etatstr = "debut"
			return etatstr
		} else {
			fmt.Println("error at getstatus function")
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

func GetNbrStatus(db *sql.DB, status string, postID int) int {
	count := 0
	query := "SELECT COUNT(*) FROM LikesDislikes WHERE post_id = ? AND " + status + " = true"
	err := db.QueryRow(query, postID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			fmt.Println("Erreur lors de l'exécution de la requête dans GetNbrStatus:", err)
			return 0
		}
	}
	return count
}
