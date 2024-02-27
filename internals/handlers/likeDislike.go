package handlers

import (
	"database/sql"
	"time"
)

type LikeDislike struct {
	PostID          int
	UserID          int
	LikeDislikeType string
	CreationDate    time.Time
}

func UpdateLikeDislike(db *sql.DB, likeDislike LikeDislike) error {
	// Check if the user has already liked or disliked the post
	var existingLikeDislikeType string
	err := db.QueryRow("SELECT like_dislike_type FROM LikesDislikes WHERE post_id = ? AND user_id = ?", likeDislike.PostID, likeDislike.UserID).Scan(&existingLikeDislikeType)
	// User has already interacted with this post
	// User is trying to like a post they already liked or dislike a post they already disliked
	// User is changing their like to a dislike or vice versa
	// User has not interacted with this post before, proceed with insertion
	shouldReturn, returnValue := newFunction(err, existingLikeDislikeType, likeDislike, db)
	if shouldReturn {
		return returnValue
	}
	return err // Error occurred during query
}

func newFunction(err error, existingLikeDislikeType string, likeDislike LikeDislike, db *sql.DB) (bool, error) {
	if err == nil {
		switch {
		case existingLikeDislikeType == likeDislike.LikeDislikeType:
			// User is trying to like or dislike the post/comment they already liked or disliked, delete the existing entry
			_, err = db.Exec("DELETE FROM LikesDislikes WHERE post_id = ? AND user_id = ?", likeDislike.PostID, likeDislike.UserID)
			if err != nil {
				return true, err
			}
		default:
			// Update the like or dislike for the post
			_, err := db.Exec("UPDATE LikesDislikes SET like_dislike_type = ?, creation_date = ? WHERE post_id = ? AND user_id = ?", likeDislike.LikeDislikeType, time.Now(), likeDislike.PostID, likeDislike.UserID)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	} else if err == sql.ErrNoRows {
		// User has not interacted with this post or comment before, proceed with insertion
		query := `
            INSERT INTO LikesDislikes (post_id, user_id, like_dislike_type, creation_date)
            VALUES (?, ?, ?, ?)
        `
		_, err := db.Exec(query, likeDislike.PostID, likeDislike.UserID, likeDislike.LikeDislikeType, likeDislike.CreationDate)
		if err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}
