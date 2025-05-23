package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

type AllData struct {
	Posts []PostWithUser
}

type PostWithUser struct {
	Post database.Post
	User database.User
}

func CreateCookie(w http.ResponseWriter) http.Cookie {
	Tokens, _ := uuid.NewV4()
	now := time.Now()
	expires := now.Add(time.Hour * 1)
	cookie := http.Cookie{
		Name:     "ForumCookie",
		Value:    Tokens.String(),
		Expires:  expires,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	return cookie
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("ForumCookie")
	if err != nil {
		return ""
	}
	return (cookie.Value)
}

func CookieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path == "/" && r.Method == "GET" {
		found := false
		ActualCookie := GetCookieHandler(w, r)
		if ActualCookie == "" {
			AllData, err := getAll(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			donnees := Data{
				Mylike:  0,
				Mypost:  0,
				Alldata: AllData,
			}
			utils.FileService("home.html", w, donnees)
			return
		}
		datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, data := range datas {
			u := data.(*database.Session)
			if u.Cookie_value == ActualCookie {
				found = true
				CurrentUser := database.User{}
				query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
				err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				CurrentUser.Firstname = utils.Trimname(CurrentUser.Firstname)
				CurrentUser.Lastname = utils.Trimname(CurrentUser.Lastname)
				AllData, err := getAll(r)
				if err != nil {
					fmt.Println(err)
					return
				}
				mylike, _ := TotalLikesByUserID(db, CurrentUser.UserID)
				mypost, _ := TotalPostByUserID(db, CurrentUser.UserID)
				notifs, errnotif := NotifTo(w, r, CurrentUser.UserID, db)
				if errnotif != nil {
					utils.FileService("error.html", w, Err[500])
					return
				}
				donnees := Data{
					Status:        "logout",
					ActualUser:    CurrentUser,
					Isconnected:   true,
					Mylike:        mylike,
					Mypost:        mypost,
					Alldata:       AllData,
					LenNotif:      notifs.LenNotif,
					Notifications: notifs.Notifs,
				}
				utils.FileService("home.html", w, donnees)
				return
			}
		}
		if !found {
			AllData, err := getAll(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			donnees := Data{
				Mylike:      0,
				Mypost:      0,
				Isconnected: false,
				Alldata:     AllData,
			}
			utils.FileService("home.html", w, donnees)
			return
		}
	} else {
		Handleredirect(w, r, db)
	}
}

func FindUserbyCookie(w http.ResponseWriter, r *http.Request, db *sql.DB, CurrentUser database.User) (bool, database.User, bool) {
	ActualCookie := GetCookieHandler(w, r)
	datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
	if err != nil {
		fmt.Println(err.Error())
		return false, database.User{}, true
	}
	Found := false
	for _, data := range datas {
		u := data.(*database.Session)
		if u.Cookie_value == ActualCookie {
			CurrentUser = database.User{}
			query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
			err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
			if err != nil {
				fmt.Println(err.Error())
				return false, database.User{}, true
			}
			Found = true
			break
		}
	}

	if r.URL.Path == "/filter" && len(datas) == 0 {
		fmt.Println("filter sans compte")
		CurrentUser.UserID = 0
	}

	return Found, CurrentUser, false
}

func TotalLikesByUserID(db *sql.DB, userID int) (int, error) {
	var totalLikes int
	query := `
        SELECT SUM(likes_count) AS total_likes
        FROM (
            SELECT COUNT(*) AS likes_count
            FROM LikesDislikes
            WHERE post_id IN (
                SELECT post_id
                FROM Posts
                WHERE user_id = ?
            ) AND liked = true
            GROUP BY post_id
        ) AS likes_per_post
    `
	err := db.QueryRow(query, userID).Scan(&totalLikes)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'exécution de la requête: %v", err)
	}
	return totalLikes, nil
}

func TotalPostByUserID(db *sql.DB, userID int) (int, error) {
	var totalPost int
	query := `SELECT COUNT(*) FROM Posts WHERE user_id = ?`
	err := db.QueryRow(query, userID).Scan(&totalPost)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'exécution de la requête: %v", err)
	}
	return totalPost, nil
}

func CountCommentsByPostID(db *sql.DB, postID int) (int, error) {
	var commentCount int
	query := "SELECT COUNT(*) FROM Comments WHERE post_id = ?"
	err := db.QueryRow(query, postID).Scan(&commentCount)
	if err != nil {
		return 0, err
	}
	return commentCount, nil
}

func Handleredirect(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var CurrentUser database.User
	var Found bool
	var shouldReturn bool
	Found, CurrentUser, shouldReturn = FindUserbyCookie(w, r, db, CurrentUser)
	if shouldReturn {
		return
	}
	switch {
	case r.URL.Path == "/" && r.Method == "POST" :
		if Found {
            PostHandler(w, r, CurrentUser)
        } else {
            utils.FileService("login.html", w, nil)
        }
	case r.URL.Path == "/filter" && r.Method == "POST":
		FilterHandler(w, r, CurrentUser)
	case r.URL.Path == "/remove" && r.Method == "GET":
		if Found {
            Removepost(w, r, CurrentUser)
        } else {
            utils.FileService("login.html", w, nil)
        }
	case r.URL.Path == "/removecomment" && r.Method == "GET":
		if Found {
            Removecomment(w, r, CurrentUser)
        } else {
            utils.FileService("login.html", w, nil)
        }
	case r.URL.Path == "/edit":
		if Found {
            Editpost(w, r, CurrentUser)
        } else {
            utils.FileService("login.html", w, nil)
        }
	case r.URL.Path == "/editcomment":
		if Found {
            Editcomment(w, r, CurrentUser)
        } else {
            utils.FileService("login.html", w, nil)
        }
	default:
		fmt.Println("au cas ou")
		w.WriteHeader(500)
		utils.FileService("error.html", w, Err[500])
		return
	}
}
