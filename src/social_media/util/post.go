package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"social_media/models"
	"time"
)

// GetPosts gets the first 10 top-level posts after `offset`.
func GetPosts(offset int, db *sql.DB) (postsJson string, err error) {

	err = errors.New("could not get posts")
	query := `
	SELECT id, post_datetime, username, post_message, COALESCE(media, '') 
	FROM posts 
	WHERE (id_replying_to IS NULL or id_replying_to = '')  
	ORDER BY post_datetime DESC 
	OFFSET $1 
	LIMIT 10;`

	if rows, err := db.Query(query, offset); err == nil {
		defer rows.Close()

		var posts []models.Post
		var post models.Post

		// get all posts from rows
		for rows.Next() {
			if err = rows.Scan(&post.Id, &post.UnixTime, &post.Username, &post.Message, &post.Media); err == nil {
				post.ReplyingTo = ""
				if post.NumReplies, err = getPostNumReplies(post.Id, db); err == nil {
					posts = append(posts, post)
				}
			}
		}
		// nothing went wrong with iterating over the rows
		if err = rows.Err(); err == nil {
			// there should be posts
			if len(posts) > 0 {
				postArr, _ := json.Marshal(posts)
				return string(postArr) + "\n", nil
			}
		}
	}
	return postsJson, err
}

// WritePost takes a request `r`, parses and writes the post to the posts table in `db` and returns post id and status code.
func WritePost(r *http.Request, db *sql.DB) (id string, statusCode int) {
	// if something goes wrong it's the request's fault
	statusCode = http.StatusBadRequest
	var post models.Post
	// read post body
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		json.Unmarshal(body, &post)
		// get sessionid cookie
		if sessionId, err := r.Cookie("sessionid"); err == nil {
			// get user name
			if post.Username, err = GetUserName(db, sessionId.Value); err == nil {
				post.Id = uniqueID()
				post.UnixTime = time.Now().Unix()
				post.Message = html.EscapeString(post.Message)

				query := `
					INSERT INTO posts (id, post_datetime, username, post_message, id_replying_to) 
					VALUES ($1, $2, $3, $4, $5);
					`
				// write post to db
				if _, err = db.Exec(query, post.Id, post.UnixTime, post.Username, post.Message, post.ReplyingTo); err == nil {
					statusCode = http.StatusOK
					id = post.Id
				} else {
					// we couldn't write the post
					statusCode = http.StatusInternalServerError
				}
			} else {
				// we couldn't get the username
				statusCode = http.StatusInternalServerError
			}
		}
	}
	return id, statusCode
}

// GetPost gets a post by `id` from the posts table in `db`.
func GetPost(id string, db *sql.DB) (string, error) {
	var post models.Post
	err := errors.New("couldn't get the post by id")
	postJson := []byte("")
	query := `
	SELECT id, post_datetime, username, post_message, 
	COALESCE(id_replying_to, ''), COALESCE(media, '')
	FROM posts 
	WHERE id = $1;`

	if rows, err := db.Query(query, id); err == nil {
		defer rows.Close()
		rows.Next()
		if err = rows.Scan(&post.Id, &post.UnixTime, &post.Username, &post.Message, &post.ReplyingTo, &post.Media); err == nil {
			// get the number of replies
			if post.NumReplies, err = getPostNumReplies(post.Id, db); err == nil {
				postJson, err = json.Marshal(post)
				return string(postJson), err
			}
		}
	}
	return string(postJson), err
}

// Get20Replies gets the first 20 top-level replies after `offset` for a post with `id`.
func Get20Replies(id string, offset int, db *sql.DB) (string, error) {
	query := `
	SELECT id, post_datetime, username, post_message, COALESCE(media, '')
	FROM posts 
	WHERE id_replying_to = $1 
	ORDER BY post_datetime 
	OFFSET $2 
	LIMIT 20;`

	var posts []models.Post
	var post models.Post
	postArr := ""
	err := errors.New("could not get replies ")
	if rows, err := db.Query(query, id, offset); err == nil {
		defer rows.Close()
		// get all posts from rows
		for rows.Next() {
			if err = rows.Scan(&post.Id, &post.UnixTime, &post.Username, &post.Message, &post.Media); err == nil {
				post.ReplyingTo = id
				if post.NumReplies, err = getPostNumReplies(post.Id, db); err == nil {
					posts = append(posts, post)
				}
			}
		}
		// nothing went wrong with iterating over the rows
		if err = rows.Err(); err == nil {
			// there should be posts
			if len(posts) > 0 {
				postArrBytes, err := json.Marshal(posts)
				return string(postArrBytes), err
			}
		}
	}
	return postArr, err
}

// DeletePost deletes the post with `id`, if it belongs to user identified by `sessionid` cookie.
func DeletePost(id string, r *http.Request, db *sql.DB) error {
	sessionId, err := r.Cookie("sessionid")
	if err != nil {
		return err
	}

	username, err := GetUserName(db, sessionId.Value)
	if err != nil {
		return err
	}

	query := `
	DELETE FROM posts
	WHERE id = $1
	AND username = $2;`
	result, err := db.Exec(query, id, username)
	if err != nil {
		return err
	}

	// if the post doesn't belong to the user, rowsAffected will be 0
	if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected < 1 {
		return errors.New("rows affected error")
	}

	// delete any media associated with the post
	os.RemoveAll("./media/" + username + "/" + id)

	return nil
}

// getPostNumReplies returns the number of replies for post with id `id`.
func getPostNumReplies(id string, db *sql.DB) (numReplies int, err error) {
	query := `
	SELECT COUNT(1) 
	FROM posts 
	WHERE id_replying_to = $1;`
	numRepliesRow, err := db.Query(query, id)
	if err != nil {
		return -1, err
	}

	numRepliesRow.Next()
	err = numRepliesRow.Scan(&numReplies)
	if err != nil {
		return -1, err
	}

	numRepliesRow.Close()

	return numReplies, nil
}

// uniqueID returns a probably unique id to identify a post with.
func uniqueID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		// if there's an error, good luck charlie
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
