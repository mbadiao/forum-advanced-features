package handlers

import (
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request, CurrentUser database.User) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		checkboxfilter := r.Form["Category"]
		if len(checkboxfilter) == 0 {
			fmt.Println("Empty checkbox filter")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if !utils.CheckCategory(checkboxfilter) {
			fmt.Println("Bad request: Invalid category")
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}

		Isconnected := utils.Isconnected(CurrentUser)

		categorypost, createdlikedpost, foundAll := utils.SplitFilter(checkboxfilter)

		query, noquery := utils.QueryFilter(categorypost, createdlikedpost, foundAll, Isconnected, CurrentUser)
		

		if noquery == "err" {
			data := Data{
				Page: "signin",
			}
			fmt.Println("filtre sans login")
			utils.FileService("login.html", w, data)
			return
		}


		AllData, err1 := getAllFilter(w, r, query, categorypost)
		if len(AllData.Posts) == 0 {
			utils.FileService("error.html", w, Err[0])
			return
		}
		if err1 != nil {
			w.WriteHeader(400)
			utils.FileService("error.html", w, Err[400])
			return
		}
		var donnees Data
		mylike,_:=TotalLikesByUserID(db,CurrentUser.UserID)
		mypost,_:=TotalPostByUserID(db,CurrentUser.UserID)
		CurrentUser.Firstname=utils.Trimname(CurrentUser.Firstname)
		CurrentUser.Lastname=utils.Trimname(CurrentUser.Lastname)
		if Isconnected {
			
			donnees = Data{
				Status:      "logout",
				ActualUser:  CurrentUser,
				Isconnected: true,
				Mylike: mylike,
				Mypost: mypost,
				Alldata:     AllData,
			}
		} else {
			donnees = Data{
				Status:      "login",
				Isconnected: false,
				Mylike: mylike,
				Mypost: mypost,
				Alldata:     AllData,
			}
		}

		utils.FileService("home.html", w, donnees)
		return
	} else {
		fmt.Println("filter method different de POST")
		w.WriteHeader(405)
		utils.FileService("error.html", w, Err[405])
		return
	}
}
