package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo PUT
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID dal path
	userID := ps.ByName("user_id")

	// Estraggo l'id dell'utente richiesto
	toFollowID := ps.ByName("followed_user_id")

	// Controllo se lo user vuole seguire se stesso
	if userID == toFollowID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check di autorizzazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Controllo se l'utente richiesto esiste (da mettere 404)
	_, err := rt.db.GetUsername(toFollowID)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user: error executing GetUsername db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'utente che si vuole seguire è già seguito
	res, err := rt.db.IsFollowed(userID, toFollowID)
	if res {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("ban-user: Error occurred while checking if the user to be banned is already banned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'utente che si vuole seguire è stato bannato
	res, err = rt.db.IsBan(userID, toFollowID)
	if res {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user: error executing isBan query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se l'utente che si vuole seguire ha bannato l'utente richiedente
	res, err = rt.db.IsBan(toFollowID, userID)
	if res {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user: error executing isBan query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Procedo a seguire l'utente richiesto
	err = rt.db.FollowUserDB(userID, toFollowID)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user: error executing FollowUser db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
