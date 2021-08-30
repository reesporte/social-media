package util

import (
	"database/sql"
	"io/ioutil"
	"mime/multipart"
	"os"
	"unicode"
)

// UploadMedia stores the user uploaded `file` and adds the new resource url to the appropriate record in db.
func UploadMedia(id string, filename string, file multipart.File, size int64, sessionId string, db *sql.DB) (err error) {
	defer file.Close()

	username, err := GetUserName(db, sessionId)
	if err != nil {
		return err
	}

	fileBytes := make([]byte, size)
	if bytes, err := file.Read(fileBytes); err != nil || bytes < int(size) {
		return err
	}

	dir := "./media/" + username + "/" + id
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0777)

	// remove illegal chars
	fname := make([]rune, 0)
	for _, char := range filename {
		if unicode.IsLetter(char) || char == '.' {
			fname = append(fname, char)
		}
	}

	filename = string(fname)

	fpath := "/media/" + username + "/" + id + "/" + filename
	if err = ioutil.WriteFile("."+fpath, fileBytes, 0777); err != nil {
		return err
	}

	// add the url to the db
	query := `
		UPDATE users
		SET media = $1
		WHERE username = $2;
	`

	if id != "profile" {
		query = `
		UPDATE posts
		SET media = $1
		WHERE id = $2;
		`
		username = id
	}

	_, err = db.Exec(query, fpath, username)
	if err != nil {
		return err
	}

	return nil
}
