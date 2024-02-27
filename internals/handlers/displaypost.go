package handlers

import (
	"database/sql"
	"fmt"
	"forum/internals/database"
	"forum/internals/utils"
	"net/http"
)

func getPosts(db *sql.DB, query string, args ...interface{}) ([]database.Post, error) {
	var posts []database.Post
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post database.Post
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.PhotoURL, &post.Content, &post.CreationDate); err != nil {
			return nil, err
		}
		categories, err := GetPostCategories(db, post.PostID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		post.FormatedDate = utils.FormatTimeAgo(post.CreationDate)

		liked := GetStatus(db, "liked", post.PostID, post.UserID)
		post.StatusLiked = liked
		disliked := GetStatus(db, "disliked", post.PostID, post.UserID)
		post.StatusDisliked = disliked

		post.Nbrlike = GetNbrStatus(db, "liked", post.PostID)
		post.Nbrdislike = GetNbrStatus(db, "disliked", post.PostID)
		nbrcomments, _ := CountCommentsByPostID(db, post.PostID)
		post.Nbrcomments = nbrcomments

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func getUser(db *sql.DB, userID int) (database.User, error) {
	var user database.User
	query := "SELECT user_id, username, firstname, lastname, email, password_hash, registration_date FROM Users WHERE user_id = ?"
	err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.PasswordHash, &user.RegistrationDate)
	return user, err
}

func getPostsWithUser(db *sql.DB, query string, args ...interface{}) ([]PostWithUser, error) {
	posts, err := getPosts(db, query, args...)
	if err != nil {
		return nil, err
	}

	var postsWithUser []PostWithUser
	for _, post := range posts {
		user, err := getUser(db, post.UserID)
		if err != nil {
			fmt.Println("Error fetching user for post:", err)
			continue
		}
		postsWithUser = append(postsWithUser, PostWithUser{
			Post: post,
			User: user,
		})
	}
	return postsWithUser, nil
}

func getAll(r *http.Request) (AllData, error) {
	query := "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts ORDER BY creation_date DESC"
	postsWithUser, err := getPostsWithUser(db, query)
	if err != nil {
		return AllData{}, err
	}
	DATA := AllData{
		Posts: postsWithUser,
	}
	return DATA, nil
}

func getAllFilter(w http.ResponseWriter, r *http.Request, query string, categorypost []string) (AllData, error) {
	var categoryInterfaces []interface{}
	for _, cat := range categorypost {
		categoryInterfaces = append(categoryInterfaces, cat)
	}
	postsWithUser, err := getPostsWithUser(db, query, categoryInterfaces...)
	if err != nil {
		return AllData{}, err
	}
	DATA := AllData{
		Posts: postsWithUser,
	}
	return DATA, nil
}

func getAllcomment(w http.ResponseWriter, r *http.Request, query string) (AllData, error) {
	postsWithUser, err := getPostsWithUser(db, query)
	if err != nil {
		return AllData{}, err
	}
	DATA := AllData{
		Posts: postsWithUser,
	}
	return DATA, nil
}