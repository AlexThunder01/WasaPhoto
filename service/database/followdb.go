package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetFollowed(userID string) ([]User, error) {
	// Creo un array vuoto dove mettero i followed
	var followed []User

	// Prendo le istanze dal db di tutte i followed
	rows, err := db.c.Query("SELECT followed_id FROM followers WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	// Metto il followed_id di ogni istanza nell'array followed
	for rows.Next() {
		var followedID User
		err = rows.Scan(&followedID)
		if err != nil {
			return nil, err
		}
		followed = append(followed, followedID)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return followed, nil

}

func (db *appdbimpl) GetFollowersID(UserID string) ([]string, error) {
	// Creo un array vuoto dove mettero i followersID
	var followersID []string

	// Prendo le istanze dal db di tutti i followers
	rows, err := db.c.Query("SELECT user_id FROM followers WHERE followed_id = ?", UserID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	// Metto il followerID di ogni istanza nell'array followers
	for rows.Next() {
		var followerID string
		err = rows.Scan(&followerID)
		if err != nil {
			return nil, err
		}
		followersID = append(followersID, followerID)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return followersID, nil
}

func (db *appdbimpl) GetFollowedsID(userID string) ([]string, error) {
	// Creo un array vuoto dove mettero i followed
	var followedsID []string

	// Prendo le istanze dal db di tutte i followed
	rows, err := db.c.Query("SELECT followed_id FROM followers WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	// Metto il followed_id di ogni istanza nell'array followed
	for rows.Next() {
		var followedID string
		err = rows.Scan(&followedID)
		if err != nil {
			return nil, err
		}
		followedsID = append(followedsID, followedID)
	}

	err = rows.Err()
	// ?
	if err != nil {
		return nil, err
	}

	return followedsID, nil

}

func (db *appdbimpl) FollowUserDB(userID string, toFollowID string) error {
	_, err := db.c.Exec(`INSERT INTO followers(user_id, followed_id) VALUES(?,?)`, userID, toFollowID)

	return err
}

func (db *appdbimpl) UnFollowUserDB(userID string, toUnfollowID string) error {
	_, err := db.c.Exec(`DELETE FROM followers WHERE user_id = ? AND followed_id = ?`, userID, toUnfollowID)

	return err
}

func (db *appdbimpl) IsFollowed(UserID string, toFollowID string) (bool, error) {

	query := "SELECT followed_id FROM followers WHERE user_id = ? AND followed_id = ?"

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
