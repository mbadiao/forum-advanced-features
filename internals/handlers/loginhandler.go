package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Data struct {
	ActualUser    database.User
	Page          string
	Badpassword   string
	Messagelg     string
	Messagesg     string
	Status        string
	Isconnected   bool
	Mylike        int
	Mypost        int
	Alldata       AllData
	LenNotif      int
	Notifications []database.Notifications
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := database.CreateTable()
	if r.URL.Path == "/login" {
		if r.Method == "GET" {
			daTa := Data{
				Page: "signin",
			}
			found := false
			datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
			if err != nil {
				fmt.Println("data")
				fmt.Println(err.Error())
				return
			}
			ActualCookie := GetCookieHandler(w, r)
			for _, data := range datas {
				s := data.(*database.Session)
				if s.Cookie_value == ActualCookie {
					found = true
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			}
			if !found {
				utils.FileService("login.html", w, daTa)
				return
			}
		} else if r.Method == "POST" {
			var data Data
			if (Empty(r.FormValue("login-name")) && Empty(r.FormValue("login-password")) &&
				Empty(r.FormValue("firstname")) && Empty(r.FormValue("lastname")) && Empty(r.FormValue("username")) &&
				Empty(r.FormValue("signup-email")) && Empty(r.FormValue("signup-password"))) || (Empty(r.FormValue("login-name")) && !Empty(r.FormValue("login-password"))) ||
				(!Empty(r.FormValue("login-name")) && Empty(r.FormValue("login-password"))) {
				data = Data{
					Page:      "signin",
					Messagelg: "All fields must be completed",
				}
				w.WriteHeader(400)
				utils.FileService("login.html", w, data)
				return
			} else if !utils.IsValidEmail(r.FormValue("signup-email")) && (Empty(r.FormValue("login-name")) && Empty(r.FormValue("login-password"))) {
				fmt.Println(r.FormValue("signup-email"))
				data = Data{
					Page:      "signup",
					Messagesg: "Email must be valid",
				}
				w.WriteHeader(400)
				utils.FileService("login.html", w, data)
				return
			} else if !utils.IsValidPassword(r.FormValue("signup-password")) && (Empty(r.FormValue("login-name")) && Empty(r.FormValue("login-password"))) {
				fmt.Println(len(r.FormValue("signup-password")))
				data = Data{
					Page:      "signup",
					Messagesg: "Password must contain at least 5 characters.",
				}
				w.WriteHeader(400)
				utils.FileService("login.html", w, data)
				return
			} else if !Empty(r.FormValue("login-name")) && !Empty(r.FormValue("login-password")) {
				var (
					id             int
					passwordhashed string
				)
				err := db.QueryRow("SELECT user_id, password_hash FROM Users WHERE email = ? OR username = ?", strings.ToLower(r.FormValue("login-name")), strings.ToLower(r.FormValue("login-name"))).Scan(&id, &passwordhashed)
				if err != nil {
					data := Data{
						Page:      "signin",
						Messagelg: "Invalid Email or Username",
					}
					w.WriteHeader(400)
					utils.FileService("login.html", w, data)
					return
				}
				err = bcrypt.CompareHashAndPassword([]byte(passwordhashed), []byte(r.FormValue("login-password")))
				if err != nil {
					data := Data{
						Page:        "signin",
						Badpassword: "Invalid Password",
					}
					w.WriteHeader(400)
					utils.FileService("login.html", w, data)
					return
				}
				if CheckAndModifySession(db, id, w, r) {
					return
				}
			} else {
				firstname, err1 := IsEmpty(r.FormValue("firstname"))
				lastname, err2 := IsEmpty(r.FormValue("lastname"))
				username, err3 := IsEmpty(r.FormValue("username"))
				email, err4 := IsEmpty(r.FormValue("signup-email"))
				password, err5 := IsEmpty(r.FormValue("signup-password"))
				if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
					data1 := Data{
						Page:      "signup",
						Messagesg: "all fields must be completed",
					}
					w.WriteHeader(400)
					utils.FileService("login.html", w, data1)
					return
				} else {
					if !utils.IsAlphaSpace(firstname) || !utils.IsAlphaSpace(lastname) {
						data1 := Data{
							Page:      "signup",
							Messagesg: "use alphanumeric characters between 2 and 15",
						}
						w.WriteHeader(400)
						utils.FileService("login.html", w, data1)
						return
					}
					var nbremail, nbrusername int
					err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email=?", strings.ToLower(email)).Scan(&nbremail)
					err1 := db.QueryRow("SELECT COUNT(*) FROM Users WHERE username=?", strings.ToLower(username)).Scan(&nbrusername)
					if err1 != nil || err != nil {
						fmt.Println("database error")
						w.WriteHeader(500)
						utils.FileService("error.html", w, Err[500])
						return
					}
					if nbremail > 0 {
						fmt.Println("email already used")
						data := Data{
							Page:      "signup",
							Messagesg: "Email already used",
						}
						w.WriteHeader(405)
						utils.FileService("login.html", w, data)
						return
					}
					if nbrusername > 0 {
						fmt.Println("Username already used")
						data := Data{
							Page:      "signup",
							Messagesg: "Username already used",
						}
						w.WriteHeader(405)
						utils.FileService("login.html", w, data)
						return
					}
					hashedpassword, errr := bcrypt.GenerateFromPassword([]byte(password), 5)
					if errr != nil {
						fmt.Println("failed to generate password")
						w.WriteHeader(500)
						utils.FileService("error.html", w, Err[500])
						return
					}
					database.Insert(db, "Users", "(username, firstname, lastname, email, password_hash)", strings.ToLower(username), firstname, lastname, strings.ToLower(email), string(hashedpassword))
				}
				data := Data{
					Page: "signin",
				}
				utils.FileService("login.html", w, data)
				return
			}
		} else {
			w.WriteHeader(405)
			utils.FileService("error.html", w, Err[405])
			return
		}
	} else {
		fmt.Println("404")
		w.WriteHeader(404)
		utils.FileService("error.html", w, Err[404])
		return
	}
}

func CheckAndModifySession(db *sql.DB, id int, w http.ResponseWriter, r *http.Request) bool {
	found := false
	datas, err := database.Scan(db, "SELECT * FROM SESSIONS ", &database.Session{})
	if err != nil {
		fmt.Println("data")
		fmt.Println(err.Error())
		return true
	}
	for _, data := range datas {
		u := data.(*database.Session)
		if u.UserID == id {
			found = true
			db.Exec(`DELETE FROM Sessions WHERE user_id =` + strconv.Itoa(id))
			database.Insert(db, "Sessions", "(user_id, cookie_value)", id, CreateCookie(w).Value)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return true
		}
	}
	if !found {
		database.Insert(db, "Sessions", "(user_id, cookie_value)", id, CreateCookie(w).Value)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return true
	}
	return false
}

func IsEmpty(str string) (string, error) {
	if strings.TrimSpace(str) == "" {
		return "", fmt.Errorf("all fields must be completed")
	}
	return strings.TrimSpace(str), nil
}

func Empty(str string) bool {
	if strings.TrimSpace(str) == "" {
		return true
	} else {
		return false
	}
}
