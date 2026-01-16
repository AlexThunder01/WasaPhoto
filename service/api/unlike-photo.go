package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo l'id dell'autore della photo
	photoOwnerID := ps.ByName("user_id")

	// Estraggo la photoID dal path
	photoID := ps.ByName("photo_id")

	// Estraggo il likeID dal path
	likeID := ps.ByName("like_id")

	// Estraggo l'id dell'utente che vuole togliere il like
	requestingUserID := extractAuthToken(r.Header.Get("Authorization"))

	// Controllo se l'utente richiedente è loggato
	if requestingUserID == "" {
		ctx.Logger.Error("like-photo: requestingUser is not logged")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Controllo se l'autore della photo ha bannato l'utente
	res, err := rt.db.IsBan(photoOwnerID, requestingUserID)
	if res {
		ctx.Logger.WithError(err).Error("unlikephoto: Photo author banned the user")
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("unlikephoto: error executing IsLiked")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'utente richiedente ha messo like alla photo
	res, err = rt.db.IsLiked(requestingUserID, photoID)
	if !res {
		ctx.Logger.WithError(err).Error("unlikephoto: user doesn't like the photo")
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("unlikephoto: error executing IsLiked")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'id dell'utente richiedente è uguale all'id del like (Devono essere la stessa cosa)
	if requestingUserID != likeID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Procedo a togliere il like
	err = rt.db.UnlikePhotoDB(requestingUserID, photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikephoto: error executing UnLikePhoto db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
