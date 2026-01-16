package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo dal body lo username e lo metto nell'oggetto user
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ctx.Logger.WithError(err).Error("do-login: error extracting bodydata")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se lo username è valido
	if !isValidUsername(user.Username) {
		http.Error(w, "the username entered is invalid", http.StatusBadRequest)
		return
	}

	// Controllo se l'utente esiste già
	res, err := rt.db.GetUserID(user.Username)

	if err != nil {
		ctx.Logger.WithError(err).Error("dologin: error executing GetUserID query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if res != -1 {
		// Se l'utente esiste già
		user.UserID = res
	} else {
		// Se l'utente non esiste lo creo
		userID, err := rt.db.CreateUser(user.Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("doLogin: error executing CreateUser query")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userIDconv, _ := strconv.ParseInt(userID, 10, 64)
		user.UserID = userIDconv
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}

}
