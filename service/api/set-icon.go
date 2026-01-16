package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setUserIcon(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Leggo il corpo della richiesta
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in the body request", http.StatusBadRequest)
		return
	}

	// Converto il corpo della richiesta in una stringa
	newIconID := string(body)

	// Controllo se il corpo della richiesta è un icon_id (intero compreso tra 5 e 20)
	err = isValidIconID(newIconID)
	if err != nil {
		ctx.Logger.WithError(err).Error("set-icon: The new iconID is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Aggiorno l'icon nel db
	err = rt.db.UpdateUserIcon(userID, newIconID)
	if err != nil {
		ctx.Logger.WithError(err).Error("set-icon: error executing UpdateUserIcon query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Invio la risposta http 204
	w.WriteHeader(http.StatusNoContent)

}

func isValidIconID(iconID string) error {

	iconIDint, err := strconv.Atoi(iconID)
	if err != nil {
		return err
	}

	if iconIDint < 5 || iconIDint > 20 {
		err = errors.New("l'iconID non è nel range giusto")
		return err
	}

	return nil
}
