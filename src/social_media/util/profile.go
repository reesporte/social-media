package util

import (
	"database/sql"
	"social_media/models"
)

// UpdateBio updates the profile `p` with sessionid `sessionId` in the users table in `db`.
func UpdateBio(p models.Profile, sessionId string, db *sql.DB) error {
	query := `
	UPDATE users
	SET bio = $1
	WHERE sessionid = $2;
	`

	result, err := db.Exec(query, p.Bio, sessionId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return err
	}

	return nil
}
