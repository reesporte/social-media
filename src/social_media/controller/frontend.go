package controller

import (
	"database/sql"
	"net/http"
	"social_media/util"

	_ "github.com/lib/pq"
)

// Home handles requests to '/' by serving the 'frontend/views/home.html' file.
func Home(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.LoggedIn(r, db) {
			http.ServeFile(w, r, "./frontend/views/home.html")
			return
		}
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
}

// Static serves static content to requests to the '/static' endpoint by serving files from the 'frontend/static' directory.
func Static() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend"+r.URL.Path)
	}
}

// Robots serves the robots.txt file to requests to the '/robots.txt' endpoint.
func Robots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./robots.txt")
	}
}
