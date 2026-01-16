package api

import (
	"encoding/json"
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Controllo se il metodo è di tipo PUT
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID
	userID := ps.ByName("user_id")

	// Verifica autenticazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Estraggo il nuovo username dal body della richiesta
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controllo se lo username è valido
	if !isValidUsername(user.Username) {
		http.Error(w, "The username must contain only characters of the alphabet or numbers", http.StatusBadRequest)
		return
	}

	// Controllo se lo username esiste gia'
	res, err := rt.db.GetUserID(user.Username)
	if res != -1 {
		http.Error(w, "The username already exists", http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error checking if username already exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Aggiorno lo username nel db
	err = rt.db.UpdateUsername(userID, user.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName: error executing updateUsername query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Invio la risposta http 204
	w.WriteHeader(http.StatusNoContent)

}
