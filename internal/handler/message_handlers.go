package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"twitter-feed/internal/model"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type MessageHandler struct {
	DB *gorm.DB
}

func (m *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	flusher := w.(http.Flusher)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")

	db, err := m.DB.DB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := db.Query("EXPERIMENTAL CHANGEFEED FOR twitterdb.messages;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notifyDisconnected := r.Context().Done()
	defer rows.Close()
	for rows.Next() {
		select {
		case <-notifyDisconnected:
			return
		default:
			var topic string
			var key string
			var value string

			err = rows.Scan(&topic, &key, &value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			message, err := extractMessage(value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			flusher.Flush()
		}
	}
}

func (m *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message model.Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if message.Username == "" || message.Message == "" {
		http.Error(w, "Fields username and message are mandatory", http.StatusBadRequest)
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := []byte(fmt.Sprintf("address-%s", r.RemoteAddr))

	err = pushToKafka(key, messageBytes, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message added"))
}

func extractMessage(value string) (model.Message, error) {
	var payload struct {
		After model.Message `json:"after"`
	}
	err := json.Unmarshal([]byte(value), &payload)
	if err != nil {
		return model.Message{}, err
	}

	return payload.After, err
}

func pushToKafka(key []byte, message []byte, context context.Context) error {
	url := os.Getenv("KAFKA_URL")
	topic := os.Getenv("KAFKA_TOPIC")
	writer := &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	msg := kafka.Message{
		Key:   key,
		Value: message,
	}

	return writer.WriteMessages(context, msg)
}
