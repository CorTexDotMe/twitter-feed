package handler

import (
	"encoding/json"
	"net/http"
	"twitter-feed/internal/model"

	"gorm.io/gorm"
)

type MessageHandler struct {
	DB *gorm.DB
}

func (m *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("New message")
}

func (m *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message model.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode("Message added")
}
