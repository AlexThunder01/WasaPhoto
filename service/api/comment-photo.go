package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo l'id del proprietario della photo dal path
	PhotoOwnerID := ps.ByName("user_id")

	// Estraggo la photoID dal path
	photoID := ps.ByName("photo_id")

	// Estraggo l'id dell'utente richiedente
	requestingUserID := extractAuthToken(r.Header.Get("Authorization"))

	// Controllo se l'utente richiedente è loggato
	if requestingUserID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Leggo il corpo della richiesta
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in the body request", http.StatusBadRequest)
		return
	}

	// Converto i dati del corpo in una stringa
	text := string(body)

	// Controllo se il testo è valido
	if !isValidText(text) {
		http.Error(w, "The comment must be less than 50 characters", http.StatusBadRequest)
		return
	}

	// Controllo se l'utente richiedente è stato bannato dal proprietario della photo
	res, err := rt.db.IsBan(PhotoOwnerID, requestingUserID)
	if res {
		ctx.Logger.WithError(err).Error("commentPhoto: Photo owner banned the requesting user")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error running IsBan")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Metto il commento alla photo
	commentID, err := rt.db.CommentPhotoDB(requestingUserID, photoID, text)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in add comment to the db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Converto lo requestingUserID in un int64
	userIDconv, err := strconv.ParseInt(requestingUserID, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error on converting userID in int64")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Converto la photoID in un int64
	photoIDconv, err := strconv.ParseInt(photoID, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error on converting phptoID in int64")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Creo un commento
	var comment Comment
	comment.AuthorID = userIDconv
	comment.CommentID = commentID
	comment.PhotoID = photoIDconv
	comment.Text = text
	comment.Datetime = GetDateTime()

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error encoding comment in json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func isValidText(text string) bool {
	return len(text) >= 1 && len(text) <= 50
}
