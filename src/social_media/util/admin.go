package util

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"social_media/models"
)

// CreateUser adds a new profile `p` to the users table.
func CreateUser(p models.Profile, db *sql.DB) (password string, err error) {
	if err = userExists(p, db); err != nil {
		return "", err
	}

	query := `
		INSERT INTO users(username, pw, join_date, bio, media, user_group)
		VALUES ($1, $2, extract(epoch from now()), 'new user!', '', $3);
	`

	password = generatePassword()
	pw := hashAndSalt(password)

	group := "not-admin"
	if p.Admin {
		group = "admin"
	}

	_, err = db.Exec(query, p.Username, pw, group)

	return password, err
}

// DeleteUser removes a profile `p` from the users table.
func DeleteUser(p models.Profile, db *sql.DB) (err error) {
	query := `
	DELETE FROM users
	WHERE username = $1;
	`

	_, err = db.Exec(query, p.Username)

	return err
}

// ResetUserPassword resets the password of profile `p` in the users table.
func ResetUserPassword(p models.Profile, db *sql.DB) (password string, err error) {
	query := `
	UPDATE users
	SET pw = $1
	WHERE username = $2;
	`
	password = generatePassword()
	pw := hashAndSalt(password)

	_, err = db.Exec(query, pw, p.Username)

	return password, err
}

// generatePassword generates a random password.
func generatePassword() string {
	adjs := []string{"Incredible", "Wonderful", "Delightful", "Lovely"}
	nouns := []string{"Beep", "Boop", "Bop", "Blork", "Blonk"}

	return choice(adjs) + choice(nouns) + fmt.Sprint(rand.Intn(100))
}

// userExists checks if the user already exists
func userExists(p models.Profile, db *sql.DB) (err error) {
	query := `
		SELECT count(1)
		FROM users
		WHERE username = $1;
	`

	exist := 0
	if existRow, err := db.Query(query, p.Username); err == nil {
		existRow.Next()
		if err = existRow.Scan(&exist); err != nil {
			return err
		}
		existRow.Close()
	} else {
		return err
	}

	if exist > 0 {
		return errors.New("user already exists")
	}

	return nil
}

// choice returns a random item from `list`.
func choice(list []string) string {
	return list[rand.Intn(len(list))]
}
