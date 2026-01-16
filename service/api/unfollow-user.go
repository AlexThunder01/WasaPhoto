package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID dal path
	userID := ps.ByName("user_id")

	// Estraggo l'id dell'utente richiesto
	toUnfollowID := ps.ByName("followed_user_id")

	// Verifica autenticazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Controllo se l'utente vuole unfolloware se stesso
	if userID == toUnfollowID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// ##CONTROLLARE SE AGGIUNGERE CHECK SE UTENTE VUOLE SEGUIRE ALTRI UTENTI CHE NON ESISTONO O BANNATI

	// Procedo a togliere il follow
	err := rt.db.UnFollowUserDB(userID, toUnfollowID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollow-user: error executing unFollowUser db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
