package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"twitter-feed/internal/database"
	"twitter-feed/internal/model"

	kafka "github.com/segmentio/kafka-go"
)

func PullFromKafka() {
	db := database.NewPostgresConnection()

	url := os.Getenv("KAFKA_URL")
	groupId := os.Getenv("KAFKA_GROUP")
	topic := os.Getenv("KAFKA_TOPIC")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{url},
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	fmt.Println("Start consuming...")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		feedMessage := model.Message{}
		err = json.Unmarshal(msg.Value, &feedMessage)
		if err != nil {
			log.Fatal(err)
		}
		db.Create(&feedMessage)

		fmt.Printf("Message created with id: %d\n", feedMessage.ID)
	}

}
