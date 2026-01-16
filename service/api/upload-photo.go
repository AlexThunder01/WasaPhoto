package api

import (
	"bytes"
	"encoding/json"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

var photoFolder = filepath.Join("/tmp", "media")

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo se il metodo è di tipo POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Estraggo lo userID dal path
	userID := ps.ByName("user_id")

	// Verifica autenticazione
	check := checkUserToken(userID, extractAuthToken(r.Header.Get("Authorization")))
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Leggo il contenuto del body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Riaggiorno il body letto
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Controllo se l'immagine è nel formato jpg o png
	err = checkImageFormat(r.Body, data)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: photo data must be jpg or png")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Riaggiorno il body letto
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Prendo il path delle photo dell'utente
	photoPath, err := getUserPhotoFolder(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: photo path error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ottengo l'istante attuale
	oraAttuale := time.Now()

	// Formatto la data e l'ora
	dataOra := oraAttuale.Format("2006-01-02 15:04:05")

	// Creo nel db un'istanza di photo e ne ritorno l'id
	photoID, err := rt.db.MakePhoto(userID, dataOra)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error in creating the photodb")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Controllo se la cartella, e se non esiste la creo
	if err := os.MkdirAll(photoPath, 0777); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("upload-photo: error creating directory")
		return
	}

	// Converto photoID in una string
	photoIDstr := strconv.FormatInt(photoID, 10)

	file, err := os.Create(photoPath + "/" + photoIDstr)
	if err != nil {
		http.Error(w, "Errore nella creazione del file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Copia direttamente i dati dal corpo della richiesta al file usando un buffer
	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Errore nel salvataggio del file", http.StatusInternalServerError)
		return
	}

	// Converto lo userID in un int64
	userIDint64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("upload-photo: error on converting UserID in int64")
		return
	}

	// Creo una struttura di tipo photo
	photo := Photo{
		AuthorID: userIDint64,
		Datetime: dataOra,
		PhotoID:  photoID,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("upload-photo: error encoding Photo")
		return
	}

}

func checkImageFormat(body io.ReadCloser, data []byte) error {
	_, err := png.Decode(body)
	if err != nil {
		body = io.NopCloser(bytes.NewBuffer(data))
		_, err := jpeg.Decode(body)
		if err != nil {
			return err
		}
	}
	return nil
}

func getUserPhotoFolder(user_id string) (string, error) {

	photoPath := filepath.Join(photoFolder, user_id, "photos")

	return photoPath, nil
}
