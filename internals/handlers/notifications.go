package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strings"
)

type NotifData struct {
	LenNotif int
	Notifs   []database.Notifications
}

func WhoLike(w http.ResponseWriter, r *http.Request, db *sql.DB, post_id int, liked string, disliked string, comuser string) {
	notif := database.Notifications{}
	actualcookie := GetCookieHandler(w, r)
	var UserId int
	title := TakePostText(post_id)
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =?", actualcookie).Scan(&UserId)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	u := database.User{}
	err = db.QueryRow("SELECT username FROM Users WHERE user_id =?", UserId).Scan(&u.Username)
	if err != nil {
		fmt.Println(err.Error())
		utils.FileService("error.html", w, Err[500])
		return
	}
	if comuser != "" {
		err = db.QueryRow("SELECT user_id FROM Posts WHERE post_id =?", post_id).Scan(&notif.UserID)
		if err != nil {
			fmt.Println(err.Error())
			utils.FileService("error.html", w, Err[500])
			return
		}
		if UserId != notif.UserID {
			notif.Message = fmt.Sprintln("", "comment your post about: ", title, "...")
		} else {
			return
		}
	}
	if liked != "" || disliked != "" {
		var likestatus bool
		var dislikestatus bool
		err = db.QueryRow("SELECT liked, disliked FROM LikesDislikes WHERE post_id =? and user_id=?", post_id, UserId).Scan(&likestatus, &dislikestatus)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var existingNotificationID int
		err = db.QueryRow("SELECT notification_id FROM Notifications WHERE post_id = ? AND user_id = ?", post_id, UserId).Scan(&existingNotificationID)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err.Error())
			return
		}
		if existingNotificationID != 0 {
			return
		}
		err = db.QueryRow("SELECT user_id FROM Posts WHERE post_id =?", post_id).Scan(&notif.UserID)
		if err != nil {
			fmt.Println("test")
			fmt.Println(err.Error())
			utils.FileService("error.html", w, Err[500])
			return
		}
		if UserId != notif.UserID {
			if liked != "" && likestatus {
				notif.Message = fmt.Sprintln("", liked, "your post about: ", title, "...")
			} else if disliked != "" && dislikestatus {
				notif.Message = fmt.Sprintln("", disliked, "your  post about: ", title, "...")
			}
		} else {
			return
		}
		if !likestatus && !dislikestatus {
			return
		}
	}
	fmt.Printf("liked %v disliked %v comment %v", liked, disliked, comuser)
	if strings.TrimSpace(liked) != "" || strings.TrimSpace(disliked) != "" || strings.TrimSpace(comuser) != "" {
		notif.Wholike = u.Username
		notif.PostID = post_id
		database.Insert(db, "Notifications", "(user_id, message, post_id,username)", notif.UserID, notif.Message, notif.PostID, notif.Wholike)
	}
}
func NotifTo(w http.ResponseWriter, r *http.Request, user_id int, db *sql.DB) (NotifData, error) {
	var Data NotifData
	actualcookie := GetCookieHandler(w, r)
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE cookie_value =? and user_id=?", actualcookie, user_id).Scan(&user_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	var notifs []database.Notifications
	rows, getnotiferr := db.Query("SELECT * FROM Notifications WHERE user_id =? ORDER BY creation_date DESC", user_id)
	if getnotiferr != nil {
		fmt.Println(getnotiferr.Error())
		return Data, nil
	}
	defer rows.Close()
	for rows.Next() {
		var notif database.Notifications
		err = rows.Scan(&notif.NotificationID, &notif.UserID, &notif.Message, &notif.PostID, &notif.Wholike, &notif.Read, &notif.CreationDate)
		if err != nil {
			fmt.Println(err.Error())
		}
		notif.Formatdate = utils.FormatTimeAgo(notif.CreationDate)
		notifs = append(notifs, notif)
	}
	Data.LenNotif = len(notifs)
	Data.Notifs = notifs
	return Data, nil
}
func TakePostText(post_id int) string {
	var title string
	errfrompost := db.QueryRow("SELECT title FROM Posts WHERE post_id=?", post_id).Scan(&title)
	if errfrompost != nil {
		return ""
	}
	return title
}

func AlreadyLiked() {

}
