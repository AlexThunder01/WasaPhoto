package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) LikePhotoDB(UserID string, photoID string) error {
	_, err := db.c.Exec("INSERT INTO likes(user_id, photo_id) VALUES(?,?)", UserID, photoID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnlikePhotoDB(userID string, photoID string) error {
	_, err := db.c.Exec(`DELETE FROM likes WHERE user_id = ? AND photo_id = ?`, userID, photoID)

	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) IsLiked(UserID string, photoID string) (bool, error) {

	query := "SELECT user_id FROM likes WHERE user_id = ? AND photo_id = ?"

	var result string
	err := db.c.QueryRow(query, UserID, photoID).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
