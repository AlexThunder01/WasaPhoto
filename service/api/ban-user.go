package api

import (
	"net/http"
	"strconv"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo PUT
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID dal path
	userID := ps.ByName("user_id")

	// Estraggo l'id dell'utente richiesto
	toBanID := ps.ByName("ban_user_id")

	// Verifica autenticazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Controllo se l'utente che si vuole bannare è se stesso
	if userID == toBanID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controllo se l'utente da bannare è gia stato bannato
	res, err := rt.db.IsBan(userID, toBanID)
	if res {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: Error occurred while checking if the user to be banned is already banned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Banno l'utente
	err = rt.db.BanUserDB(userID, toBanID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: error in banning the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Tolgo il follow all'utente bannato se veniva seguito
	err = rt.db.UnFollowUserDB(userID, toBanID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: error to unfollow the banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Viceversa a sopra
	err = rt.db.UnFollowUserDB(toBanID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: error to unfollow the user by the banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// rimuovo i likes e i commenti da tutte le photo dell 'utente bannato e viceversa
	err = rt.rmvLikesAndComments(userID, toBanID, w, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: error in rmvLikesAndComments function 1")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.rmvLikesAndComments(toBanID, userID, w, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: error in rmvLikesAndComments function 2")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) rmvLikesAndComments(userID string, toBanID string, w http.ResponseWriter, ctx reqcontext.RequestContext) error {

	// Prendo le photo del profilo dell'utente bannato
	banUserPhotos, err := rt.db.GetProfilePhotos(toBanID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: error to get the profile's photos of the banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	for _, photo := range banUserPhotos {

		// estraggo l'id della photo
		photoID := strconv.FormatInt(photo.PhotoID, 10)

		// Tolgo il like alla photo se esiste
		err = rt.db.UnlikePhotoDB(userID, photoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("ban-user: error to unlike the profile's photos of the banned user")
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		// Tolgo i commenti alla photo se esistono
		err = rt.db.RemoveAllCommentDB(userID, photoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("ban-user: error to uncooment the profile's photos of the banned user")
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

	}
	return nil
}
