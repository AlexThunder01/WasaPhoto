package api

import (
	"encoding/json"
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo user_id dell'utente richiesto dai parametri
	userID := ps.ByName("user_id")

	// Estraggo l'id dell'utente richiedente
	loggedUserID := extractAuthToken(r.Header.Get("Authorization"))

	// Controllo se l'utente richiedente è loggato
	if loggedUserID == "" {
		ctx.Logger.Error("get-user-icon: The user who made the request is not logged in")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Faccio richiesta al db per estrarre lo username dell'utente richiesto
	username, err := rt.db.GetUsername(userID)
	if username == "" {
		ctx.Logger.WithError(err).Error("get-username: GetUsername user doesn't exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("get-username: GetUsername error getting username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
