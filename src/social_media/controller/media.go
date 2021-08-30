package controller

import (
	"database/sql"
	"net/http"
	"social_media/handler"
	"social_media/util"
)

// Media handles requests to the '/media' endpoint by routing to the appropriate handler based on request method.
func Media(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.LoggedIn(r, db) {
			switch r.Method {
			case "GET":
				http.ServeFile(w, r, "."+r.URL.Path)
			case "POST":
				handler.PostMedia(w, r, db)
			}
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
