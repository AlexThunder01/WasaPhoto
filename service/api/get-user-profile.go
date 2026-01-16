package api

import (
	"encoding/json"
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo l'id dell'utente richiedente
	userID := extractAuthToken(r.Header.Get("Authorization"))

	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Estraggo l'id dell'utente richiesto
	requestedUserID := ps.ByName("user_id")

	// Controllo se l'utente richiesto esiste ed estraggo il suo username
	username, err := rt.db.GetUsername(requestedUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-user-profile: error running GetUsername")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if username == "" {
		ctx.Logger.Error("get-user-profile: requestedUser doesn't exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Controllo se l'utente richiesto ha bannato l'utente richiedente
	result, err := rt.db.IsBan(requestedUserID, userID)
	if result {
		ctx.Logger.WithError(err).Error("get-user-profile: error, The owner of the requested profile has banned the requesting user")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("get-user-profile: error checking whether the requesting user has been banned by the requested user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prendo l'icona del profilo
	iconID, err := rt.db.GetUserIconID(requestedUserID)
	if iconID == "" {
		ctx.Logger.WithError(err).Error("get-user-icon: GetUserIconID user doesn't exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("get-user-icon: GetUserIconID error getting userIconID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Creo un oggetto di tipo Profile
	var profile Profile

	// Lo converto nel tipo ProfileDB
	profileDB := toProfileDB(profile)

	// Assegno i seguenti attributi in modo da creare un profilo parziale (caso ban)
	profileDB.IconID = iconID
	profileDB.Owner = username

	// Controllo se l'utente richiedente ha bannato l'utente richiesto
	result, err = rt.db.IsBan(userID, requestedUserID)
	if result {
		// Se l'utente richiedente ha bannato l'utente richiesto restituisco il profilo parziale
		w.WriteHeader(http.StatusPartialContent)
	} else if err != nil {
		ctx.Logger.WithError(err).Error("get-user-profile: error checking whether the requesting user has been banned by the requested user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		// Caso in cui l'utente richiesto non è stato bannato

		// Prendo i followers
		followers, err := rt.db.GetFollowersID(requestedUserID)
		if err != nil {
			ctx.Logger.WithError(err).Error("get-user-profile: error executing getFollowers query")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Prendo i followed
		followeds, err := rt.db.GetFollowedsID(requestedUserID)
		if err != nil {
			ctx.Logger.WithError(err).Error("get-user-profile: error executing getFollowed query")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Prendo le photo
		Photos, err := rt.db.GetProfilePhotos(requestedUserID)
		if err != nil {
			ctx.Logger.WithError(err).Error("get-user-profile: error executing getFollowed query")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Converto followers e followed in int 64
		followersInt64, err := convertSliceTo64(followers)
		if err != nil {
			ctx.Logger.WithError(err).Error("get-user-profile: error converting followersID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		followedsInt64, err := convertSliceTo64(followeds)
		if err != nil {
			ctx.Logger.WithError(err).Error("get-user-profile: error converting followedsID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Assegno alla struttura profile i seguenti
		profileDB.FollowedsID = followedsInt64
		profileDB.FollowersID = followersInt64
		profileDB.Photos = Photos

		w.WriteHeader(http.StatusOK)
	}

	// Converto il profilo in un json
	profileJSON, err := json.Marshal(profileDB)
	if err != nil {
		http.Error(w, "Errore durante la creazione della risposta JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(profileJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-user-profile: error on writing json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
