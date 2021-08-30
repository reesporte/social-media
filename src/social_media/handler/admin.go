package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"social_media/models"
	"social_media/util"
)

// GetAdmin serves the admin page.
func GetAdmin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/views/admin.html")
}

// CreateUser parses an '/admin/create' POST request, and passes the data to the util.CreateUser function, returning the appropriate status code to the client.
func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	statusCode := http.StatusInternalServerError
	pw := ""
	var newProfile models.Profile
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err = json.Unmarshal(body, &newProfile); err == nil {
			if pw, err = util.CreateUser(newProfile, db); err == nil {
				statusCode = http.StatusOK
			} else {
				// couln't make the new user
				statusCode = http.StatusConflict
			}
		}
	} else {
		// couldn't parse the request
		statusCode = http.StatusBadRequest
	}

	w.WriteHeader(statusCode)
	w.Write([]byte("{\"pw\": \"" + pw + "\"}"))
}

// DeleteUser parses an '/admin/' DELETE request, and passes the data to the util.DeleteUser function, returning the appropriate status code to the client.
func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	statusCode := http.StatusInternalServerError
	var profile models.Profile
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err = json.Unmarshal(body, &profile); err == nil {
			if err = util.DeleteUser(profile, db); err == nil {
				statusCode = http.StatusOK
			}
		}
	} else {
		statusCode = http.StatusBadRequest
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

// ResetUserPassword parses an '/admin/' PUT request, and passes the data to the util.ResetUserPassword function, returning the appropriate status code to the client.
func ResetUserPassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	statusCode := http.StatusInternalServerError
	pw := ""
	var profile models.Profile
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err = json.Unmarshal(body, &profile); err == nil {
			if pw, err = util.ResetUserPassword(profile, db); err == nil {
				statusCode = http.StatusOK
			}
		}
	} else {
		// couldn't read the body
		statusCode = http.StatusBadRequest
	}

	w.WriteHeader(statusCode)
	w.Write([]byte("{\"pw\": \"" + pw + "\"}"))
}
