package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"social_media/models"
	"social_media/util"
)

// Login authenticates the username and password and sets the sessionID if appropriate, returning the appropriate status code to the client.
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	statusCode := http.StatusBadRequest

	if body, err := ioutil.ReadAll(r.Body); err == nil {
		var l models.Login
		json.Unmarshal(body, &l)

		statusCode = http.StatusForbidden
		if util.AuthUserAndPass(l.Username, l.Password, db) {
			statusCode = setSessionID(w, l.Username, db)
		}
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

// Logout logs out a logged in user.
func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	c, _ := r.Cookie("sessionid")
	if user, err := util.GetUserName(db, c.Value); err == nil {
		if !util.SetSessionId("", user, db) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	// kindly ask browser to delete cookie
	c.MaxAge = -1
	http.SetCookie(w, c)
}

// setSessionID sets the session ID cookie, and adds it to the db, returning the appropriate status code to the caller.
func setSessionID(w http.ResponseWriter, username string, db *sql.DB) int {
	session := util.SessionId()
	name := "sessionid"
	raw := name + "=" + session
	cookie := http.Cookie{
		Name:       name,
		Value:      session,
		Path:       "/",
		RawExpires: "",
		MaxAge:     86400,
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
		HttpOnly:   true,
		Raw:        raw,
		Unparsed:   []string{raw}}

	if util.SetSessionId(session, username, db) {
		http.SetCookie(w, &cookie)
		return http.StatusOK
	}
	return http.StatusUnauthorized
}
