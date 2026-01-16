package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo PUT
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo l'id dell'autore della photo
	photoOwnerID := ps.ByName("user_id")

	// Estraggo la photoID dal path
	photoID := ps.ByName("photo_id")

	// Estraggo il likeID dal path
	likeID := ps.ByName("like_id")

	// Estraggo l'id dell'utente che vuole mettere like
	requestingUserID := extractAuthToken(r.Header.Get("Authorization"))

	// Controllo se l'autore della photo ha bannato l'utente
	bool, err := rt.db.IsBan(photoOwnerID, requestingUserID)
	if bool {
		ctx.Logger.WithError(err).Error("like-photo: Photo Owner banned the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("like-photo: error executing IsBan")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'id dell'utente richiedente è uguale all'id del like (Devono essere la stessa cosa)
	if requestingUserID != likeID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controllo se l'utente ha gia messo like alla photo
	res, err := rt.db.IsLiked(requestingUserID, photoID)
	if res {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("like-photo: error executing isLiked query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Metto il like alla photo
	err = rt.db.LikePhotoDB(requestingUserID, photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("likePhoto: error executing LikePhotoDB query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
