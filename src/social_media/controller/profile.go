package controller

import (
	"database/sql"
	"net/http"
	"social_media/handler"
	"social_media/util"
)

// Profile handles requests to the '/profile' endpoint by routing to the appropriate handler based on request method and URL path.
func Profile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.LoggedIn(r, db) {
			switch r.Method {
			case "GET":
				splitPath := util.SplitPath(r.URL.Path)
				if len(splitPath) > 1 && splitPath[1] != "" {
					handler.GetProfile(w, r, splitPath[1], db)
				}
			case "PUT":
				splitPath := util.SplitPath(r.URL.Path)
				if len(splitPath) > 1 && splitPath[1] != "" {
					switch splitPath[1] {
					case "password":
						handler.UpdatePassword(w, r, db)
					case "bio":
						handler.UpdateBio(w, r, db)
					}
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
			return
		}
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
}
