package controller

import (
	"database/sql"
	"net/http"
	"social_media/handler"
	"social_media/util"
)

// Admin handles requests to the '/admin' endpoint. It routes requests according to request method.
func Admin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.LoggedIn(r, db) && util.IsAdmin(r, db) {
			switch r.Method {
			case "GET":
				handler.GetAdmin(w, r)
			case "POST":
				splitPath := util.SplitPath(r.URL.Path)
				if len(splitPath) > 1 && splitPath[1] != "" {
					switch splitPath[1] {
					case "create":
						handler.CreateUser(w, r, db)
					}
				}
			case "DELETE":
				handler.DeleteUser(w, r, db)
			case "PUT":
				handler.ResetUserPassword(w, r, db)
			}
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
