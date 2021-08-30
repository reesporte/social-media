package handler

import (
	"database/sql"
	"net/http"
	"social_media/util"
)

// PostMedia parses a '/media/' POST request, and passes the data to the util.UploadMedia function, returning the appropriate statusCode to the client.
func PostMedia(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	statusCode := http.StatusInternalServerError

	if sessionId, err := r.Cookie("sessionid"); err == nil {
		if err = r.ParseMultipartForm(10 << 20); err == nil {
			if file, header, err := r.FormFile("media"); err == nil {
				id := r.FormValue("postId")
				if id == "" {
					id = "profile"
				}
				if err = util.UploadMedia(id, header.Filename, file, header.Size, sessionId.Value, db); err == nil {
					statusCode = http.StatusOK
				}
			} else {
				// we couldn't read the file
				statusCode = http.StatusBadRequest
			}
		} else {
			// payload too large
			statusCode = 413
		}
	} else {
		// don't have a sessionid cookie!
		statusCode = http.StatusBadRequest
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}
