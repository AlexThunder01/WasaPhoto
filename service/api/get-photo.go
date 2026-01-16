package api

import (
	"net/http"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// (401)

	userID := ps.ByName("user_id")
	photoID := ps.ByName("photo_id")

	http.ServeFile(w, r, photoFolder+"/"+userID+"/photos/"+photoID)
}
