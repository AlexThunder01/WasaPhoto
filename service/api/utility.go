package api

import (
	"regexp"
	"strconv"
	"time"
	"wasaphoto/service/database"
)

func toProfileDB(profile Profile) database.Profile {

	profileDB := database.Profile{
		Owner: profile.Owner,
	}

	return profileDB
}

func toPhoto(photosDB []database.Photo) []Photo {

	var photosConv []Photo

	for _, photodb := range photosDB {
		photoConv := Photo{
			AuthorID: photodb.AuthorID,
			PhotoID:  photodb.PhotoID,
			Datetime: photodb.Datetime,
		}

		photosConv = append(photosConv, photoConv)
	}

	return photosConv
}

func toLike(LikesDB []database.Like) []Like {

	var LikesConv []Like

	for _, Likedb := range LikesDB {
		LikeConv := Like{
			UserID: Likedb.UserID,
		}

		LikesConv = append(LikesConv, LikeConv)
	}

	return LikesConv
}

func toComment(CommentsDB []database.Comment) []Comment {

	var commentsConv []Comment

	for _, Commentdb := range CommentsDB {
		CommentConv := Comment{
			PhotoID:   Commentdb.PhotoID,
			AuthorID:  Commentdb.AuthorID,
			Text:      Commentdb.Text,
			CommentID: Commentdb.CommentID,
			Datetime:  Commentdb.Datetime,
		}

		commentsConv = append(commentsConv, CommentConv)
	}

	return commentsConv
}

func isValidUsername(username string) bool {
	regex := regexp.MustCompile(`^.*?[a-zA-Z0-9]+.*?$`)
	return regex.MatchString(username) && len(username) >= 3 && len(username) <= 16
}

func convertSliceTo64(slice []string) ([]int64, error) {
	var sliceStr []int64
	for i := 0; i < len(slice); i++ {
		// Converto un elemento str dello slice in int64
		el, err := strconv.ParseInt(slice[i], 10, 64)

		if err != nil {
			return nil, err
		}

		sliceStr = append(sliceStr, el)
	}

	return sliceStr, nil
}

func GetDateTime() string {
	currentTime := time.Now()
	dateFormat := "2006-01-02 15:04:05"
	dateTime := currentTime.Format(dateFormat)
	return dateTime
}
