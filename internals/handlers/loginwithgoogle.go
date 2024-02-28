package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	clientID     = "652043065465-e28rhk8h2sakg8navg377s34p12c5p3c.apps.googleusercontent.com"
	clientSecret = "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"
	redirectURI  = "http://localhost:8080/callback"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIN   int    `json:"expires_in"`
}

type UserInfoResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func HandleLogGoogle(w http.ResponseWriter, r *http.Request) {
	google_auth_URL := "https://accounts.google.com/o/oauth2/auth"
	scope := "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email"
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=code", google_auth_URL, clientID, url.QueryEscape(redirectURI), scope)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func exchangeCodeForToken(code string) (*TokenResponse, error) {
	token_url := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(token_url, data)
	if err != nil {
		fmt.Println("error at the exchange function")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("imposible to read the body", err)
		return nil, err
	}
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Println("impossible to unmarshal the body of resp", err)
		return nil, err
	}
	return &tokenResponse, nil
}

func getUserInfo(accessToken string) (*UserInfoResponse, error) {

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfoResponse

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")

	token, err := exchangeCodeForToken(code)

	if err != nil {
		w.WriteHeader(500)
		utils.FileService("error.html", w, Err[500])
		return
	}

	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		w.WriteHeader(500)
		utils.FileService("error.html", w, Err[500])
		return
	}

	// fmt.Println("resultat",userInfo)

	presence, er := IsUserInTheDataBase("email", strings.ToLower(userInfo.Email))
	if er != nil {
		w.WriteHeader(500)
		utils.FileService("error.html", w, Err[500])
		return
	}
	if presence {
		var id int
		err := db.QueryRow("SELECT user_id FROM Users WHERE email = ?", strings.ToLower(userInfo.Email)).Scan(&id)
		if err != nil {
			data := Data{
				Page:      "signin",
				Messagelg: "Invalid Email or Username",
			}
			w.WriteHeader(400)
			utils.FileService("login.html", w, data)
			return
		}
		shouldReturn := CheckAndModifySession(db, id, w, r)
		if shouldReturn {
			return
		}
	}
	if !presence {
		username := extractUsername(userInfo.Email)
		usernametest := username
		usernamePresence, errr := IsUserInTheDataBase("username", strings.ToLower(usernametest))
		if errr != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
		i := 0
		for usernamePresence {
			usernametest += strconv.Itoa(i)
			i++
		}
		username = usernametest

		hashedpassword, errr := bcrypt.GenerateFromPassword([]byte(userInfo.ID), 5)
		if errr != nil {
			fmt.Println("failed to generate password")
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}

		database.Insert(db, "Users", "(username, firstname, lastname, email, password_hash)", strings.ToLower(username), userInfo.FirstName, userInfo.LastName, strings.ToLower(userInfo.Email), string(hashedpassword))

		var id int
		err := db.QueryRow("SELECT user_id FROM Users WHERE email = ?", strings.ToLower(userInfo.Email)).Scan(&id)
		if err != nil {
			data := Data{
				Page:      "signin",
				Messagelg: "Invalid Email or Username",
			}
			w.WriteHeader(400)
			utils.FileService("login.html", w, data)
			return
		}
		if CheckAndModifySession(db, id, w, r) {
			return
		}
	}
}

func IsUserInTheDataBase(Column string, valeur string) (bool, error) {
	query := "SELECT COUNT(*) FROM Users WHERE " + Column + " = ?"

	// Exécuter la requête SQL
	var count int
	err := db.QueryRow(query, valeur).Scan(&count)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return false, err
	}

	// Vérifier si l'e-mail existe dans la base de données
	if count > 0 {
		fmt.Println(Column + " existe dans la base de données.")
		return true, nil
	} else {
		fmt.Println(Column + " n'existe pas dans la base de données.")
		return false, nil
	}
}

func extractUsername(str string) string {
	res := strings.Split(str, "@")

	return (res[0])
}
