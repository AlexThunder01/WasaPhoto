package database

import (
	"database/sql"
	"errors"
	"time"
)

func (db *appdbimpl) CommentPhotoDB(UserID string, photoID string, text string) (int64, error) {

	currentTime := time.Now()
	dateFormat := "2006-01-02 15:04:05"
	date := currentTime.Format(dateFormat)

	res, err := db.c.Exec("INSERT INTO comments(author_id, photo_id, text, date) VALUES(?,?,?,?)", UserID, photoID, text, date)
	if err != nil {
		return -1, err
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return commentID, nil
}

func (db *appdbimpl) RemoveAllCommentDB(userId string, photoId string) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE author_id = ? AND photo_id = ?`, userId, photoId)

	return err
}

func (db *appdbimpl) RemoveCommentDB(commentID string) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE comment_id = ?`, commentID)

	return err
}

func (db *appdbimpl) IsCommentAuthor(UserID string, CommentID string) (bool, error) {

	query := "SELECT comment_id FROM comments WHERE author_id = ? AND comment_id = ?"

	var result string
	err := db.c.QueryRow(query, UserID, CommentID).Scan(&result)
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
