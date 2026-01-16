package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) BanUserDB(UserID string, toBanID string) error {
	_, err := db.c.Exec("INSERT INTO banned_users(banner_id, banned_id) VALUES(?,?)", UserID, toBanID)

	return err
}

func (db *appdbimpl) IsBan(UserID string, toFollowID string) (bool, error) {

	query := "SELECT banner_id FROM banned_users WHERE banner_id = ? AND banned_id = ?"

	var result string
	err := db.c.QueryRow(query, UserID, toFollowID).Scan(&result)
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

func (db *appdbimpl) UnbanUserDB(userID string, toUnbanID string) error {

	_, err := db.c.Exec(`DELETE FROM banned_users WHERE banner_id = ? AND banned_id = ?`, userID, toUnbanID)

	return err
}
