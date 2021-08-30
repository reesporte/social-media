package controller

import (
	"database/sql"
	"net/http"
	"social_media/handler"
	"social_media/util"
)

// Posts handles requests to the '/posts' endpoint by routing to the appropriate handler based on request method and URL path.
func Posts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if util.LoggedIn(r, db) {
			splitPath := util.SplitPath(r.URL.Path)
			switch r.Method {
			case "GET":
				if len(splitPath) > 1 && splitPath[1] != "" {
					if splitPath[1] == "replies" {
						id := splitPath[2]
						handler.GetReplies(id, w, r, db)
					} else {
						id := splitPath[1]
						handler.GetPost(id, w, r, db)
					}
				} else {
					handler.GetPosts(w, r, db)
				}
				return
			case "POST":
				handler.WritePost(w, r, db)
				return
			case "DELETE":
				id := splitPath[1]
				handler.DeletePost(id, w, r, db)
				return
			}
		}
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
}
