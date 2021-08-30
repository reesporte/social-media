package handler

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"social_media/models"
	"social_media/util"
)

// UpdatePassword parses a '/profile/password' PUT request, and passes the data to the util.UpdatePassword function, returning the appropriate statusCode to the client.
func UpdatePassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	statusCode := http.StatusBadRequest
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		var pw models.PasswordChange
		json.Unmarshal(body, &pw)
		sessionId, _ := r.Cookie("sessionid")
		pw.SessionId = sessionId.Value

		if err = util.UpdatePassword(pw, db); err == nil {
			statusCode = http.StatusOK
		}
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

// UpdateBio parses a '/profile/bio' PUT request, and passes the data to the util.UpdateBio function, returning the appropriate statusCode to the client.
func UpdateBio(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	statusCode := http.StatusBadRequest

	if body, err := ioutil.ReadAll(r.Body); err == nil {
		var p models.Profile
		if err = json.Unmarshal(body, &p); err == nil {
			if sessionId, err := r.Cookie("sessionid"); err == nil {
				if err = util.UpdateBio(p, sessionId.Value, db); err == nil {
					statusCode = http.StatusOK
				} else {
					statusCode = http.StatusInternalServerError
				}
			}
		}
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

// GetProfile parses a '/profile/{id}' GET request, and generates a profile page based on the data, returning the appropriate statusCode to the client.
func GetProfile(w http.ResponseWriter, r *http.Request, username string, db *sql.DB) {
	query := `
		SELECT join_date, bio, COALESCE(media, '')
		FROM users
		WHERE username = $1;
	`

	var p models.Profile
	if profileRow, err := db.Query(query, username); err == nil {
		profileRow.Next()
		if err = profileRow.Scan(&p.JoinDate, &p.Bio, &p.Media); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			profileRow.Close()
			return
		}
		if p.Media == "" {
			p.Media = "/media/default.jpg"
		}
		profileRow.Close()
	}
	p.JoinDate *= 1000
	p.Username = username
	p.Owner = false
	if sessionid, err := r.Cookie("sessionid"); err == nil {
		if u, err := util.GetUserName(db, sessionid.Value); err == nil {
			if u == username {
				p.Owner = true
			}
		}
		p.Admin = util.IsAdminId(sessionid.Value, db)
	}

	tmpl := template.Must(template.ParseFiles("./frontend/views/p.html"))
	tmpl.Execute(w, p)
}
