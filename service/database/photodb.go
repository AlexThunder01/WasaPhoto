package database

import (
	"database/sql"
	"errors"
	"strconv"
)

func (db *appdbimpl) MakePhoto(author_id string, dataOra string) (int64, error) {

	res, err := db.c.Exec("INSERT INTO photos(author_id, date) VALUES(?,?)", author_id, dataOra)
	if err != nil {
		return -1, err
	}

	photoID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return photoID, nil
}

func (db *appdbimpl) DeletePhotoDB(photoID string) error {
	_, err := db.c.Exec(`DELETE FROM photos WHERE photo_id = ?`, photoID)

	return err
}

func (db *appdbimpl) IsPhotoOwner(author_id string, photoID string) (bool, error) {

	query := "SELECT author_id FROM photos WHERE author_id = ? AND photo_id = ?"

	var result string
	err := db.c.QueryRow(query, author_id, photoID).Scan(&result)
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

func (db *appdbimpl) GetProfilePhotos(userID string) ([]Photo, error) {

	// Seleziono tutte le istanze di photo relative al profilo in input in ordine cronologico inverso
	rows, err := db.c.Query("SELECT * FROM photos WHERE author_id = ? ORDER BY date DESC", userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	// leggo tutte le istanze di photo trovate nel db e le metto in photos
	var photos []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoID, &photo.AuthorID, &photo.Datetime)
		if err != nil {
			return nil, err
		}

		// Converto photoID in una stringa
		photoIDstr := strconv.FormatInt(photo.PhotoID, 10)

		// Creo l'array coi commenti della photo
		comments, err := db.GetPhotoComments(photoIDstr)
		if err != nil {
			return nil, err
		}
		photo.Comments = comments

		// Creo l'array coi like della photo
		likes, err := db.GetPhotoLikes(photoIDstr)
		if err != nil {
			return nil, err
		}
		photo.Likes = likes

		// Metto la photo nell'array photos
		photos = append(photos, photo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (db *appdbimpl) GetPhotoComments(photoID string) ([]Comment, error) {
	// Creo l'array vuoto dove metterò i commenti
	var comments []Comment

	// Seleziono tutte le istanze di commento relative alla photo in input
	rows, err := db.c.Query("SELECT * FROM comments WHERE photo_id = ?", photoID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		// Estraggo le informazioni di ogni instanza e le metto in una variabile commento
		var comment Comment
		err = rows.Scan(&comment.CommentID, &comment.PhotoID, &comment.AuthorID, &comment.Text, &comment.Datetime)
		if err != nil {
			return nil, err
		}

		// Metto il commento nell'array comments
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (db *appdbimpl) GetPhotoLikes(photoID string) ([]Like, error) {
	// Creo l'array vuoto dove metterò i like
	var likes []Like

	rows, err := db.c.Query("SELECT user_id FROM likes WHERE photo_id = ?", photoID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		// Estraggo le informazioni di ogni instanza e le metto in una variabile like
		var like Like
		err = rows.Scan(&like.UserID)
		if err != nil {
			return nil, err
		}

		// Metto il commento nell'array comments
		likes = append(likes, like)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func (db *appdbimpl) GetPhotoAuthorID(photoID string) (string, error) {
	var authorID string

	err := db.c.QueryRow("SELECT author_id FROM photos WHERE photo_id = ?", photoID).Scan(&authorID)
	if err != nil {
		return "", err
	}
	return authorID, nil
}
