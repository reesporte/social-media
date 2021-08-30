package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"social_media/util"
	"strconv"
)

// GetPosts parses a '/posts/' GET request, and passes the data to the util.GetPosts function, returning the appropriate statusCode to the client.
func GetPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	offset, _ := strconv.Atoi(r.Header.Get("Offset"))
	if posts, err := util.GetPosts(offset, db); err == nil {
		fmt.Fprintf(w, "%s", posts)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
}

// GetReplies parses a '/posts/replies/{id}' GET request, and passes the data to the util.GetReplies function, returning the appropriate statusCode to the client.
func GetReplies(id string, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	offset, _ := strconv.Atoi(r.Header.Get("Offset"))

	if posts, err := util.Get20Replies(id, offset, db); err == nil {
		fmt.Fprintf(w, "%s", posts)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
}

// GetPost parses a '/posts/{id}' GET request, and passes the data to the util.GetPost function, returning the appropriate statusCode to the client.
func GetPost(id string, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if post, err := util.GetPost(id, db); err == nil {
		fmt.Fprintf(w, "%s", post)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
}

// WritePost parses a '/posts/' POST request, and passes the data to the util.WritePost function, returning the appropriate statusCode to the client.
func WritePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var statusCode int
	postId, statusCode := util.WritePost(r, db)
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"id\": \"" + postId + "\", \"status\": \"" + http.StatusText(statusCode) + "\"}"))
}

// DeletePost parses a '/posts/{id}' DELETE request, and passes the data to the util.DeletePost function, returning the appropriate statusCode to the client.
func DeletePost(id string, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	statusCode := http.StatusInternalServerError
	if err := util.DeletePost(id, r, db); err == nil {
		statusCode = http.StatusOK
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}
