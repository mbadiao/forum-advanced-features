package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func WhoLike(w http.ResponseWriter, r *http.Request, db *sql.DB, post_id int, liked string, disliked string) {
	actualcookie := GetCookieHandler(w, r)
	var UserId int
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =?", actualcookie).Scan(&UserId)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	u := database.User{}
	err = db.QueryRow("SELECT username FROM Users WHERE user_id =?", UserId).Scan(&u.Username)
	if err != nil {
		fmt.Println(err.Error())
		utils.FileService("error.html", w, Err[500])
		return
	}
	var likestatus bool
	var dislikestatus bool
	err = db.QueryRow("SELECT liked, disliked FROM LikesDislikes WHERE post_id =? and user_id=?", post_id, UserId).Scan(&likestatus, &dislikestatus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	notif := database.Notifications{}
	err = db.QueryRow("SELECT user_id FROM Posts WHERE post_id =?", post_id).Scan(&notif.UserID)
	if err != nil {
		fmt.Println("test")
		fmt.Println(err.Error())
		utils.FileService("error.html", w, Err[500])
	}
	if liked != "" && likestatus {
		notif.Message = fmt.Sprintln(u.Username, "", liked, " post", post_id)
	} else if disliked != "" && dislikestatus {
		notif.Message = fmt.Sprintln(u.Username, "", disliked, " post", post_id)
	} else if liked != "" && !likestatus {
		notif.Message = fmt.Sprintln(u.Username, " un"+liked, " post", post_id)
	} else if disliked != "" && !dislikestatus {
		notif.Message = fmt.Sprintln(u.Username, " un"+disliked, " post", post_id)
	}
	notif.PostID = post_id
	database.Insert(db, "Notifications", "(user_id, message, post_id)", notif.UserID, notif.Message, notif.PostID)
}

type NotifData struct {
	LenNotif int
	Notifs   []database.Notifications
}

func NotifTo(w http.ResponseWriter, r *http.Request, user_id int, db *sql.DB) (NotifData, error) {
	var Data NotifData
	actualcookie := GetCookieHandler(w, r)
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =? and user_id=?", actualcookie, user_id).Scan(&user_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	var notifs []database.Notifications
	rows, getnotiferr := db.Query("SELECT * FROM Notifications WHERE user_id =?", user_id)
	if getnotiferr != nil {
		fmt.Println(getnotiferr.Error())
		return Data, nil
	}
	defer rows.Close()

	for rows.Next() {
		var notif database.Notifications
		err = rows.Scan(&notif.NotificationID, &notif.UserID, &notif.Message, &notif.PostID, &notif.Read)
		if err != nil {
			fmt.Println(err.Error())
		}
		notifs = append(notifs, notif)
	}
	Data.LenNotif = len(notifs)
	Data.Notifs = notifs
	return Data, nil
}
