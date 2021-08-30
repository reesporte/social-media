package util

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"social_media/models"

	"golang.org/x/crypto/bcrypt"
)

// you can try logging in 10 times, and then after that you're screwed
const MAX_LOGIN_ATTEMPTS int = 10

// AuthUserAndPass returns true when `username` is authenticated. The user account must not be locked and `password` must match what is stored in the database for `username`. If the passwords do not match, or the user account is locked, util.handleBadLogin is called.
func AuthUserAndPass(username string, password string, db *sql.DB) bool {
	query := `
		SELECT pw
		FROM users
		WHERE username = $1;
	`
	pw := ""
	if pwRow, err := db.Query(query, username); err == nil {
		defer pwRow.Close()
		pwRow.Next()
		if err = pwRow.Scan(&pw); err == nil {
			if comparePasswords(pw, password) && !lockedOut(username, db) {
				if err = setLoginAttempts(username, 0, db); err != nil {
					log.Println("Could not reset login attempts for user ", username, err)
				}
				return true
			} else {
				handleBadLogin(username, db)
			}
		}
	}
	return false
}

// lockedOut returns true if the account belonging to `username` is locked.
func lockedOut(username string, db *sql.DB) bool {
	query := `
		SELECT locked 
		FROM users
		WHERE username = $1;
	`
	locked := ""
	if lRow, err := db.Query(query, username); err == nil {
		defer lRow.Close()
		lRow.Next()
		if err = lRow.Scan(&locked); err == nil {
			// if user is locked out, locked == "Y"
			return locked == "Y"
		}
	}
	return true
}

// handleBadLogin increments login_attempts for `username` if login_attempts is less than MAX_LOGIN_ATTEMPTS. Otherwise, it locks the account associated with `user`.
func handleBadLogin(username string, db *sql.DB) {
	if !lockedOut(username, db) { // only do something if they're not already locked out
		query := `
		SELECT login_attempts
		FROM users
		WHERE username = $1;
	`
		attempts := 0
		if aRow, err := db.Query(query, username); err == nil {
			defer aRow.Close()
			aRow.Next()
			if err = aRow.Scan(&attempts); err == nil {
				if attempts > MAX_LOGIN_ATTEMPTS {
					if err = ModifyUserLock(username, true, db); err != nil {
						log.Println("Could not lock user ", username, err)
					}
				} else {
					if err = setLoginAttempts(username, attempts+1, db); err != nil {
						log.Println("Could not increment login attempts for user ", username, err)
					}
				}
			}
		}
	}
}

// setLoginAttempts sets the number of login attempts to `attempts` for `username` in the users table of the database.
func setLoginAttempts(username string, attempts int, db *sql.DB) (err error) {
	query := `
		UPDATE users
		SET login_attempts = $1
		WHERE username = $2;
	`
	_, err = db.Exec(query, attempts, username)
	return err
}

// ModifyUserLock sets the locked value for `username` in the users table to `locked`
func ModifyUserLock(username string, locked bool, db *sql.DB) (err error) {
	query := `
		UPDATE users
		SET locked = $1
		WHERE username = $2;
	`
	lock := "N"
	if locked {
		lock = "Y"
	}

	_, err = db.Exec(query, lock, username)
	return err
}

// SessionId generates a probably unique session ID.
func SessionId() string {
	return uniqueID()
}

// validSessionID ensures `id` exists in db, and the user account is not locked
func validSessionID(id string, db *sql.DB) bool {
	query := `
		SELECT count(1)
		FROM users
		WHERE sessionid = $1
		AND locked = 'N';
	`

	exist := 0
	if existRow, err := db.Query(query, id); err == nil {
		existRow.Next()
		if err = existRow.Scan(&exist); err != nil {
			return false
		}
		existRow.Close()
	}

	return exist == 1
}

// SetSessionId sets sessionid to `id` for `username` in the user table.
func SetSessionId(id string, username string, db *sql.DB) bool {
	query := `
	UPDATE users
	SET sessionid = $1
	WHERE username = $2;
	`
	if _, err := db.Exec(query, id, username); err != nil {
		return false
	}
	return true
}

// LoggedIn returns true if there is a valid sessionid stored in the request cookie `sessionid`
func LoggedIn(r *http.Request, db *sql.DB) bool {
	if sessionid, err := r.Cookie("sessionid"); err == nil {
		id := sessionid.Value
		return validSessionID(id, db)
	}
	return false
}

// IsAdmin returns true if there is a sessionid stored in the request cookie `sessionid` that is an adminId
func IsAdmin(r *http.Request, db *sql.DB) bool {
	if sessionid, err := r.Cookie("sessionid"); err == nil {
		id := sessionid.Value
		return IsAdminId(id, db)
	}
	return false
}

// IsAdminId makes sure sessionid is in database, is in an admin group, and account is not locked.
func IsAdminId(id string, db *sql.DB) bool {
	query := `
	SELECT COUNT(1)
	FROM users
	WHERE sessionid = $1 
	AND user_group = 'admin'
	AND locked = 'N';
	`
	if res, err := db.Query(query, id); err == nil {
		var count int
		res.Next()
		if err = res.Scan(&count); err == nil {
			return count == 1
		}
	}
	return false
}

// UpdatePassword updates a user's password if their current password and sessionid match.
func UpdatePassword(pw models.PasswordChange, db *sql.DB) error {
	username, err := GetUserName(db, pw.SessionId)
	if err != nil {
		return err
	}
	if !AuthUserAndPass(username, pw.CurrentPassword, db) {
		return err
	}

	if isInTop100(pw.NewPassword) {
		return err
	}

	query := `
	UPDATE users
	SET pw = $1
	WHERE sessionid = $2;
	`

	result, err := db.Exec(query, hashAndSalt(pw.NewPassword), pw.SessionId)
	if err != nil {
		return err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected != 1 {
		return errors.New("rows affected error")
	}

	return nil
}

// hashAndSalt hashes and salts password `pwd`.
func hashAndSalt(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash)
}

// comparePasswords checks passwords `hashedPwd` and `plainPwd` for equivalence.
func comparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

// isInTop100 checks if a password is in 2021's top 100 passwords.
func isInTop100(pw string) bool {
	top := []string{"123456", "password", "12345678", "qwerty", "123456789", "12345", "1234", "111111", "1234567", "dragon", "123123", "baseball", "abc123", "football", "monkey", "letmein", "696969", "shadow", "master", "666666", "qwertyuiop", "123321", "mustang", "1234567890", "michael", "654321", "pussy", "superman", "1qaz2wsx", "7777777", "fuckyou", "121212", "000000", "qazwsx", "123qwe", "killer", "trustno1", "jordan", "jennifer", "zxcvbnm", "asdfgh", "hunter", "buster", "soccer", "harley", "batman", "andrew", "tigger", "sunshine", "iloveyou", "fuckme", "2000", "charlie", "robert", "thomas", "hockey", "ranger", "daniel", "starwars", "klaster", "112233", "george", "asshole", "computer", "michelle", "jessica", "pepper", "1111", "zxcvbn", "555555", "11111111", "131313", "freedom", "777777", "pass", "fuck", "maggie", "159753", "aaaaaa", "ginger", "princess", "joshua", "cheese", "amanda", "summer", "love", "ashley", "6969", "nicole", "chelsea", "biteme", "matthew", "access", "yankees", "987654321", "dallas", "austin", "thunder", "taylor", "matrix", "minecraft"}

	for _, pass := range top {
		if pw == pass {
			return true
		}
	}

	return false
}
