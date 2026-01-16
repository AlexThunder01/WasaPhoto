package api

import (
	"encoding/json"
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Controllo se il metodo è di tipo GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Verifica autenticazione
	userID := extractAuthToken(r.Header.Get("Authorization"))

	// Controllo se l'utente è loggato
	if userID == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Estraggo la stringa per cui cercare gli usernames simili ad essa
	usersToSearch := r.URL.Query().Get("users")

	// Cerco gli usernames simili nel db
	users, err := rt.db.SearchUsers(usersToSearch, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("SearchUsers: error executing SearchUsers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
