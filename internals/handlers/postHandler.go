package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Category struct {
	CategoryID   int
	CategoryName string
}
var alldata AllData

var db = database.CreateTable()
var postID int64 // Declare postID outside the if statement

func PostHandler(w http.ResponseWriter, r *http.Request, userid database.User) {
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
		fmt.Println()
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
		database.Insert(db, "Posts", "(user_id, title, PhotoURL, content)", PostS.UserID, PostS.Title, PostS.PhotoURL, PostS.Content)
		err = db.QueryRow("SELECT last_insert_rowid()").Scan(&postID)
		if err != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
		CategoriesId := utils.GetCategory(CheckboxValues)
		for _, v := range CategoriesId {
			database.Insert(db, "PostCategories", "(post_id, category_id)", postID, v)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Function to retrieve categories of a post
func GetPostCategories(db *sql.DB, postID int) ([]string, error) {
	query := `
    SELECT PostCategories.post_id, Categories.name
    FROM PostCategories
    INNER JOIN Categories ON PostCategories.category_id = Categories.category_id
    WHERE PostCategories.post_id =` + strconv.Itoa(postID)
	// Call the Scan function with the PostCategory struct
	data, err := database.Scan(db, query, &database.Category{})
	if err != nil {
		return nil, err
	}
	// Extract category names from the result
	var categories []string
	for _, d := range data {
		category := d.(*database.Category)
		categories = append(categories, category.Name)
	}
	return categories, nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) string {
	var photoURL = "NoPhoto"

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		return "err400"
	}
	// Obtenez le fichier téléchargé à partir de la clé du formulaire
	file, handler, err := r.FormFile("file")
	if file != nil {
		if err != nil {
			fmt.Println("fichier" + err.Error())
			return "err400"
		}
		defer file.Close()
		if !utils.IsValidImage(file, handler) {
			return "err400"
		}
		// Vérifiez la taille du fichier
		if handler.Size > 20<<20 {
			return "err408"
		}
		dirPath := "./web/static/upload"

		// Check if the directory exists, if not create it
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				fmt.Println("répertoire" + err.Error())
				return "err500"
			}
		}

		// Créez un fichier dans le répertoire temporaire du serveur
		tempFile, err := os.CreateTemp(dirPath, "upload-*"+filepath.Ext(handler.Filename))
		if err != nil {
			fmt.Println("répertoire" + err.Error())
			return "err500"
		}
		defer tempFile.Close()
		// Copiez le contenu du fichier téléchargé dans le fichier temporaire
		_, err = io.Copy(tempFile, file)
		if err != nil {
			fmt.Println("Copiez" + err.Error())
			return "err500"
		}
		// Générez une URL pour le fichier téléchargé (par exemple, l'URL peut être le chemin relatif vers le fichier temporaire)
		photoURL = tempFile.Name()
		fmt.Println("avant", photoURL)
	}
	return photoURL
}
