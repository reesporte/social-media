package util

import (
	"database/sql"
	"fmt"
)

// These are the values for the postgres db connection.
const (
	user     = "postgres"
	dbname   = "social_media"
	password = "social_media"
	host     = "db"
	sslmode  = "disable"
)

// InitDb opens and returns a db connection using the const values.
func InitDb() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", host, user, password, dbname, sslmode))
}

// GetUserName gets a username based on sessionid.
func GetUserName(db *sql.DB, sessionId string) (username string, err error) {
	query := `
		SELECT (username)
		FROM users
		WHERE sessionid = $1;
	`
	if row, err := db.Query(query, sessionId); err == nil {
		defer row.Close()
		row.Next()
		if err = row.Scan(&username); err != nil {
			return username, err
		}
	}
	return username, nil
}
