package database

import (
	"database/sql"
)

type Table interface {
	ScanRows(rows *sql.Rows) error
}

// Pour la structure User
func (u *User) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&u.UserID, &u.Username, &u.Firstname, &u.Lastname, &u.Email, &u.PasswordHash, &u.RegistrationDate)
}

// Pour la structure Post
func (p *Post) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&p.PostID, &p.UserID, &p.Title, &p.PhotoURL, &p.Content, &p.CreationDate)
}

// Pour la structure Comment
func (c *Comment) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&c.CommentID, &c.PostID, &c.UserID, &c.Content, &c.CreationDate)
}

// Pour la structure Category
func (cat *Category) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&cat.CategoryID, &cat.Name)
}

// Pour la structure PostCategory
func (pc *PostCategory) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&pc.PostID, &pc.CategoryID)
}

// Pour la structure LikeDislike
func (ld *LikeDislike) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&ld.LikeDislikeID, &ld.PostID, &ld.CommentID, &ld.UserID, &ld.LikeDislikeType, &ld.CreationDate)
}

// Pour la structure Session
func (s *Session) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&s.SessionID, &s.UserID, &s.Cookie_value, &s.ExpirationDate)
}
