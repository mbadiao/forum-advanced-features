package handlers

import (
	"net/http"
)

type Error struct {
	Code    int
	Message string
}

var Err = map[int]Error{
	404: {
		http.StatusNotFound,
		http.StatusText(404),
	},
	500: {
		http.StatusInternalServerError,
		http.StatusText(500),
	},
	400: {
		http.StatusBadRequest,
		http.StatusText(400),
	},
	405: {
		http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed),
	},
	0: {
		0,
		"Result",
	},
	1: {
		400,
		"create a post with an image larger than 20mb is not possible",
	},
	401: {
		401,
		"user not identified, please check your information",
	},
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  []string
}

var Routes = []Route{
	{
		Path:    "/login",
		Handler: LoginHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/",
		Handler: HomeHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/filter",
		Handler: HomeHandler,
		Method:  []string{"POST"},
	},
	{
		Path:    "/comment",
		Handler: CommentHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/logout",
		Handler: LogoutHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/likedislike",
		Handler: LikeDislikeHandler,
		Method:  []string{"POST"},
	},
	{
		Path:    "/likecomment",
		Handler: LikeCommentHandler,
		Method:  []string{"POST"},
	},
	{
		Path:    "/callback",
		Handler: HandleCallback,
		Method:  []string{"GET"},
	},
	{
		Path:    "/continuewithgoogle",
		Handler: HandleLogGoogle,
		Method:  []string{"GET"},
	},
	{
		Path:    "/githubcallback",
		Handler: HandleGitHubCallback,
		Method:  []string{"GET"},
	},
	{
		Path:    "/continuewithgithub",
		Handler: HandleLogGithub,
		Method:  []string{"GET"},
	},
	{
		Path:    "/profile",
		Handler: ProfileHandler,
		Method:  []string{"POST", "GET"},
	},
	{
		Path:    "/remove",
		Handler: Removepost,
		Method:  []string{"GET"},
	},
	{
		Path:    "/edit",
		Handler: HomeHandler,
		Method:  []string{"POST", "GET"},
	},
}

var Port = ":8080"
