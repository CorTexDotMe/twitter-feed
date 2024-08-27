package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

	var messages []*model.Message
	err := m.DB.Find(&messages).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range messages {
		json.NewEncoder(w).Encode(v)
		flusher.Flush()
	}

	for {
		_, err := w.Write([]byte("Message"))
		if err != nil {
			break
		}
		flusher.Flush()

		time.Sleep(5 * time.Second)
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
	}
	key := []byte(fmt.Sprintf("address-%s", r.RemoteAddr))

	err = pushToKafka(key, messageBytes, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message added"))
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
