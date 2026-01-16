package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo l'id dell'utente richiedente
	userID := ps.ByName("user_id")

	// Verifica autenticazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Prendo gli utenti seguiti e li metto in un array
	followeds, err := rt.db.GetFollowedsID(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetMyStream: error executing GetFollowed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var Stream []Photo

	for _, followedID := range followeds {

		// Prendo le photo di ogni followed
		Photos, err := rt.db.GetProfilePhotos(followedID)
		if err != nil {
			ctx.Logger.WithError(err).Error("GetProfilePhotos: error executing GetFollowed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Trasformo Photos di tipo database.Photo nel tipo Photo
		PhotosConv := toPhoto(Photos)

		// Aggiungo i likes e i commenti ad ogni photo
		for _, Photo := range PhotosConv {
			PhotoID := strconv.FormatInt(Photo.PhotoID, 10)

			Likes, err := rt.db.GetPhotoLikes(PhotoID)
			if err != nil {
				ctx.Logger.WithError(err).Error("GetMyStream: error executing GetPhotoLikes query")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Trasformo Likes di tipo database.Photo nel tipo Photo
			LikesConv := toLike(Likes)

			Photo.Likes = LikesConv

			Comments, err := rt.db.GetPhotoComments(PhotoID)
			if err != nil {
				ctx.Logger.WithError(err).Error("GetMyStream: error executing GetPhotoComments query")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Trasformo Likes di tipo database.Photo nel tipo Photo
			CommentsConv := toComment(Comments)

			Photo.Comments = CommentsConv

			// Inserisco nell'oggetto photo lo username dell'autore

			authorIDstr := strconv.Itoa(int(Photo.AuthorID))

			username, err := rt.db.GetUsername(authorIDstr)
			if username == "" {
				ctx.Logger.WithError(err).Error("get-my-stream: GetUsername user doesn't exists")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if err != nil {
				ctx.Logger.WithError(err).Error("get-my-stream: GetUsername error getting username")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			Photo.Username = username

			Stream = append(Stream, Photo)
		}

	}

	sort.Slice(Stream, func(i, j int) bool {
		time1, _ := time.Parse(time.RFC3339, Stream[i].Datetime)
		time2, _ := time.Parse(time.RFC3339, Stream[j].Datetime)
		return time1.After(time2) // Invertito rispetto all'ordinamento crescente
	})

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(Stream)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetMyStream: error encoding json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
