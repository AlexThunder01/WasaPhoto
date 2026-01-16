/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	UpdateUsername(userID string, username string) error
	MakePhoto(profileID string, dataOra string) (int64, error)
	GetProfilePhotos(profileID string) ([]Photo, error)
	GetPhotoComments(photoID string) ([]Comment, error)
	GetPhotoLikes(photoID string) ([]Like, error)
	FollowUserDB(userID string, toFollowID string) error
	IsBan(userID string, toFollowID string) (bool, error)
	UnFollowUserDB(userID string, toUnfollowID string) error
	BanUserDB(UserID string, toBanID string) error
	UnbanUserDB(userID string, toUnbanID string) error
	GetUserID(Username string) (int64, error)
	LikePhotoDB(UserID string, photoID string) error
	IsLiked(UserID string, photoID string) (bool, error)
	UnlikePhotoDB(userID string, photoID string) error
	CommentPhotoDB(UserID string, photoID string, text string) (int64, error)
	RemoveAllCommentDB(userId string, photoId string) error
	DeletePhotoDB(photoID string) error
	CreateUser(Username string) (string, error)
	IsPhotoOwner(UserID string, photoID string) (bool, error)
	GetFollowed(userID string) ([]User, error)
	GetFollowersID(UserID string) ([]string, error)
	GetFollowedsID(userID string) ([]string, error)
	IsCommentAuthor(UserID string, CommentID string) (bool, error)
	RemoveCommentDB(commentID string) error
	GetPhotoAuthorID(photoID string) (string, error)
	GetUsername(UserID string) (string, error)
	SearchUsers(usersToSearch string, userID string) ([]User, error)
	GetUserIconID(UserID string) (string, error)
	UpdateUserIcon(userID string, iconID string) error
	IsFollowed(UserID string, toFollowID string) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func createDatabase(db *sql.DB) error {
	tables := [7]string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) NOT NULL,
			icon_id INTEGER DEFAULT 5 NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER NOT NULL,			
			photo_id INTEGER NOT NULL,			
			PRIMARY KEY (user_id, photo_id),						
			FOREIGN KEY(photo_id) REFERENCES photos (photo_id) ON DELETE CASCADE
			FOREIGN KEY(user_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			photo_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			text VARCHAR(30) NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY(photo_id) REFERENCES photos (photo_id) ON DELETE CASCADE,
			FOREIGN KEY(author_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS photos (
			photo_id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY(author_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS banned_users (
			banner_id INTEGER NOT NULL,
			banned_id INTEGER NOT NULL,
			PRIMARY KEY (banner_id, banned_id),
			FOREIGN KEY(banner_id) REFERENCES users (user_id) ON DELETE CASCADE,
			FOREIGN KEY(banned_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS followers(
			user_id INTEGER NOT NULL,
			followed_id INTEGER NOT NULL,
			PRIMARY KEY (user_id,followed_id),
			FOREIGN KEY(user_id) REFERENCES users (user_id) ON DELETE CASCADE,
			FOREIGN KEY(followed_id) REFERENCES users (user_id) ON DELETE CASCADE
			);`,
	}

	for i := 0; i < len(tables); i++ {
		table := tables[i]
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
