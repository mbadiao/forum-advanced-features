package handlers

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	query := ""
	var CurrentUser database.User
	ActualCookie := GetCookieHandler(w, r)
	datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	found := false

	for _, data := range datas {
		u := data.(*database.Session)
		if u.Cookie_value == ActualCookie {
			CurrentUser = database.User{}
			query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id=?"
			err := db.QueryRow(query, u.UserID).Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Firstname, &CurrentUser.Lastname, &CurrentUser.Email, &CurrentUser.PasswordHash, &CurrentUser.RegistrationDate)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			found = true
			break
		}
	}

	mess := ""
	if r.URL.RawQuery == "like" {
		query = "SELECT DISTINCT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date FROM Posts p JOIN LikesDislikes ld ON p.post_id = ld.post_id WHERE ld.user_id = " + strconv.Itoa(CurrentUser.UserID) + " AND (ld.liked = TRUE OR ld.disliked = TRUE) ORDER BY creation_date DESC"
		mess = "Liked or Disliked Post"
	} else if r.URL.RawQuery == "create" || r.URL.RawQuery == "" {
		query = "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts WHERE user_id =" + strconv.Itoa(CurrentUser.UserID)
		mess = "Created Post"
	} else if r.URL.RawQuery == "comment" {
		query = "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts WHERE post_id IN (SELECT post_id FROM Comments WHERE user_id = " + strconv.Itoa(CurrentUser.UserID) + ")"
		mess = "Commented Post"
	} else {
		w.WriteHeader(404)
		utils.FileService("error.html", w, Err[404])
		return
	}

	if found {
		Displayprofile(w, r, CurrentUser, query, mess)
		return
	} else {
		w.WriteHeader(404)
		utils.FileService("error.html", w, Err[404])
		return
	}
}

func Displayprofile(w http.ResponseWriter, r *http.Request, CurrentUser database.User, query, mess string) {
	Result := false
	code := 0
	AllData, err1 := getAllcomment(w, r, query)
	if len(AllData.Posts) == 0 {
		Result = true
		// utils.FileService("error.html", w, Err[0])
		// return
	}
	if err1 != nil {
		w.WriteHeader(400)
		utils.FileService("error.html", w, Err[400])
		return
	}
	var donnees Data
	mylike, _ := TotalLikesByUserID(db, CurrentUser.UserID)
	mypost, _ := TotalPostByUserID(db, CurrentUser.UserID)
	CurrentUser.Firstname = utils.Trimname(CurrentUser.Firstname)
	CurrentUser.Lastname = utils.Trimname(CurrentUser.Lastname)
	donnees = Data{
		Status:       "logout",
		ActualUser:   CurrentUser,
		Isconnected:  true,
		Mylike:       mylike,
		Results:      Result,
		Code0results: code,
		Mess0results: mess,
		Mypost:       mypost,
		Alldata:      AllData,
	}
	utils.FileService("profile.html", w, donnees)
}
