package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID del proprietario del path
	userID := ps.ByName("user_id")

	// Estraggo l'id dell'utente da sbannare
	toUnbanID := ps.ByName("ban_user_id")

	// Controllo se l'utente richiedente è proprietario del path
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Procedo a togliere il ban
	err := rt.db.UnbanUserDB(userID, toUnbanID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unban-user: error executing unbanUser db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
