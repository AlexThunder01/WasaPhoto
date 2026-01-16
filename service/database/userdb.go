package database

import (
	"database/sql"
	"errors"
	"strconv"
)

func (db *appdbimpl) CreateUser(Username string) (string, error) {
	res, err := db.c.Exec("INSERT INTO users(username) VALUES(?)", Username)
	if err != nil {
		return "", err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	userIDstr := strconv.FormatInt(userID, 10)

	return userIDstr, nil
}

// Funzione del database che aggiorna lo username
func (db *appdbimpl) UpdateUsername(userID string, username string) error {

	_, err := db.c.Exec(`UPDATE users SET username = ? WHERE user_id = ?`, username, userID)
	if err != nil {
		// Error durante l'esecuzione della query
		return err
	}
	return nil
}

func (db *appdbimpl) GetUserID(Username string) (int64, error) {
	var result int64

	err := db.c.QueryRow("SELECT user_id FROM users WHERE username = ?", Username).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		} else {
			return -1, err
		}
	} else {
		return result, nil
	}
}

func (db *appdbimpl) GetUsername(UserID string) (string, error) {
	var username string

	err := db.c.QueryRow("SELECT username FROM users WHERE user_id = ?", UserID).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", err
		}
	} else {
		return username, nil
	}

}

func (db *appdbimpl) SearchUsers(usersToSearch string, userID string) ([]User, error) {
	// Creo l'array vuoto dove metterò gli user
	var users []User

	// Seleziono tutte le istanze di gli usernames simili alla variabile usersToSearch ed al tempo stesso restituisco solo gli usernames degli utenti che non hanno bannato l'utente richiedente
	rows, err := db.c.Query("SELECT user_id, username FROM users WHERE (username LIKE ?) AND ? NOT IN (SELECT banned_id FROM banned_users WHERE banner_id = user_id)", "%"+usersToSearch+"%", userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		// Estraggo le informazioni di ogni istanza e le metto nella in un oggetto di tipo User
		var user User
		err = rows.Scan(&user.UserID, &user.Username)
		if err != nil {
			return nil, err
		}

		// Metto lo User creato nell'array degli users
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *appdbimpl) GetUserIconID(UserID string) (string, error) {
	var iconID string

	err := db.c.QueryRow("SELECT icon_id FROM users WHERE user_id = ?", UserID).Scan(&iconID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", err
		}
	} else {
		return iconID, nil
	}
}

func (db *appdbimpl) UpdateUserIcon(userID string, iconID string) error {

	_, err := db.c.Exec("UPDATE users SET icon_id = ? WHERE user_id = ?", iconID, userID)

	return err
}
