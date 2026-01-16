package api

import (
	"strings"
)

// Questa funzione ritorna true se il check è andato a buon fine
func checkUserToken(userID string, authToken string) bool {

	if authToken == "" {
		return false
	}

	if userID != authToken {
		return false
	}

	return true

}

// Funzione per estrarre il token ID dall'header di autorizzazione
func extractAuthToken(authorization string) string {

	var header = strings.Split(authorization, " ")

	if len(header) == 2 {
		return strings.Trim(header[1], " ")
	}

	return ""
}
