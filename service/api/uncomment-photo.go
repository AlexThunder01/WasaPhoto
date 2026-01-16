package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo l'id del proprietario della photo dal path
	PhotoOwnerID := ps.ByName("user_id")

	// Estraggo il commentID dal path
	commentID := ps.ByName("comment_id")

	// Estraggo l'id dell'utente richiedente
	requestingUserID := extractAuthToken(r.Header.Get("Authorization"))

	// Controllo se l'autore della photo ha bannato l'utente richiedente
	res, err := rt.db.IsBan(PhotoOwnerID, requestingUserID)
	if res {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Photo owner banned the requesting user")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error running IsBan")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'utente è il proprietario del commento
	res, err = rt.db.IsCommentAuthor(requestingUserID, commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error checking the author of the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !res {
		ctx.Logger.WithError(err).Error("uncommentPhoto: user is not the author of the comment")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Procedo a togliere il commento
	err = rt.db.RemoveCommentDB(commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error executing UncommentPhotoDB db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
