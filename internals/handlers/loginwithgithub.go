package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"io"
	"strconv"
	"strings"

	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

const (
	githubClientID     = "a87cd687fb1da946fc20"
	githubClientSecret = "d7e4f010ad7090fd217671eab0812f14dd5eff73"
	githubRedirectURI  = "http://localhost:8080/githubcallback"
)

type UserInfoGithub struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserEmail []struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}


func HandleLogGithub(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:read,user:email,read:user", githubClientID, url.QueryEscape(githubRedirectURI))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}


func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Récupération du code d'autorisation GitHub
	code := r.URL.Query().Get("code")

	// Échange du code d'autorisation contre un jeton d'accès
	token, err := exchangeCodeForTokenGithub(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Utilisation du jeton d'accès pour obtenir les informations de l'utilisateur
	userInfo, err := getUserInfoGithub(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	process(w,r,userInfo)
}


func exchangeCodeForTokenGithub(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", githubRedirectURI)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}


func getUserInfoGithub(token string) (*UserInfoGithub, error) {
	// Création de la requête pour obtenir les informations de l'utilisateur
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	// Envoi de la requête
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Lecture de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Analyse de la réponse JSON pour obtenir les informations de l'utilisateur
	var userInfo UserInfoGithub
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	req2, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}

	req2.Header.Add("Authorization", "Bearer "+token)
	resp2, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	Body, _ := io.ReadAll(resp2.Body)

	var userEmail UserEmail

	err = json.Unmarshal(Body, &userEmail)
	if err != nil {
		return nil, err
	}

	for i := range userEmail {
		if userEmail[i].Primary {
			userInfo.Email = userEmail[i].Email
			break
		}
	}

	return &userInfo, nil
}


func firstAndLastName(chaine string) (string, string, error) {

	dernierIndiceEspace := strings.LastIndex(chaine, " ")

	if dernierIndiceEspace == -1 {
		return "", "", fmt.Errorf("aucun espace trouvé dans la chaîne")
	}

	partieAvant := chaine[:dernierIndiceEspace]
	partieApres := chaine[dernierIndiceEspace+1:]

	return partieAvant, partieApres, nil
}


func process(w http.ResponseWriter, r *http.Request, allinfo *UserInfoGithub) {
	firstname, lastname, err := firstAndLastName(allinfo.Name)
	if err != nil {
		w.WriteHeader(401)
		utils.FileService("error.html", w, Err[401])
		return
	}

	presence, err := IsUserInTheDataBase("email", strings.ToLower(allinfo.Email))
	if err != nil {
		w.WriteHeader(500)
		utils.FileService("error.html", w, Err[500])
		return
	}

	if presence {
		var id int
		err := db.QueryRow("SELECT user_id FROM Users WHERE email = ?", strings.ToLower(allinfo.Email)).Scan(&id)
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

		usernamePresence, errr := IsUserInTheDataBase("username", strings.ToLower(allinfo.Login))
		if errr != nil {
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}
		i := 0
		for usernamePresence {
			allinfo.Login += strconv.Itoa(i)
			i++
		}
		hashedpassword, errr := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(allinfo.Id)), 5)
		if errr != nil {
			fmt.Println("failed to generate password")
			w.WriteHeader(500)
			utils.FileService("error.html", w, Err[500])
			return
		}

		database.Insert(db, "Users", "(username, firstname, lastname, email, password_hash)", strings.ToLower(allinfo.Login), firstname, lastname, strings.ToLower(allinfo.Email), string(hashedpassword))

		var id int
		err := db.QueryRow("SELECT user_id FROM Users WHERE email = ?", strings.ToLower(allinfo.Email)).Scan(&id)
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
