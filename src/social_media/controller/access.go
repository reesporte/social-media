package controller

import (
	"database/sql"
	"net/http"
	"social_media/handler"
	"social_media/util"
)

// Login handles requests to the '/login' endpoint. It routes requests to the correct handler based on the request method.
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !util.LoggedIn(r, db) {
			switch r.Method {
			case "GET":
				http.ServeFile(w, r, "./frontend/views/login.html")
			case "POST":
				handler.Login(w, r, db)
			}
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout routes requests from the '/logout' endpoint to the Logout handler.
func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.LoggedIn(r, db) {
			handler.Logout(w, r, db)
		}
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
}
