package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
)

func CreateMessageByUser(w http.ResponseWriter, r *http.Request) {
	var chatMessage models.ChatMessage
	err := json.NewDecoder(r.Body).Decode(&chatMessage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error en la creacion del mensaje"})
		return
	}

	db.DB.Create(&chatMessage)
	json.NewEncoder(w).Encode(&chatMessage)
}

func GetMessageByUser(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	sender := params.Get("sender")
	recipient := params.Get("recipient")

	var chatMessage []models.ChatMessage

	db.DB.Where("(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)", sender, recipient, recipient, sender).Find(&chatMessage)
	json.NewEncoder(w).Encode(&chatMessage)

}
