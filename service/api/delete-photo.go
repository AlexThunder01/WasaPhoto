package api

import (
	"net/http"
	"os"
	"path/filepath"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID dal path
	userID := ps.ByName("user_id")

	// Estraggo l'id della photo da cancellare
	photoID := ps.ByName("photo_id")

	// Verifica autenticazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Controllo se l'utente è il proprietario della photo (aggiungere 403)
	_, err := rt.db.IsPhotoOwner(userID, photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error executing isPhotoOwner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prendo il path dell'immagine della photo
	path := filepath.Join(photoFolder, userID, "/photos/"+photoID)

	// Procedo a togliere la photo dal db
	err = rt.db.DeletePhotoDB(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error executing DeletePhotoDB query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Cancello l'immagine della photo dal sistema
	err = os.Remove(path)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error to remove photo in the os")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
