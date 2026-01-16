package api

import (
	"net/http"
	"os"
	"path/filepath"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getServerIcon(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo il parametro icon_id
	iconID := ps.ByName("icon_id")

	// (401, 400)

	// Variabile che contiene il path di un icona
	iconsFolder, err := getIconsFolderPath(iconID)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-server-icon: error getting icons Folder path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Invio al client l'icona richiesta
	http.ServeFile(w, r, iconsFolder)
}

func getIconsFolderPath(iconID string) (string, error) {
	// Funzione che genera il path di una determinata icona

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	iconsFolder := filepath.Join(currentDir, "webui", "src", "assets", "icons", iconID) // Supponendo che la cartella icons sia due livelli sopra il percorso corrente
	return iconsFolder, nil
}
