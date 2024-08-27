package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"twitter-feed/internal/database"
	"twitter-feed/internal/model"

	kafka "github.com/segmentio/kafka-go"
)

func PullFromKafka() {
	db := database.NewPostgresConnection()

	url := os.Getenv("KAFKA_URL")
	groupId := os.Getenv("KAFKA_GROUP")
	topic := os.Getenv("KAFKA_TOPIC")

	sendInitMessage(url, groupId, topic)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{url},
		GroupID: groupId,
		Topic:   topic,
	})
	defer reader.Close()

	fmt.Println("Start consuming...")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		feedMessage := model.Message{}
		err = json.Unmarshal(msg.Value, &feedMessage)
		if err != nil {
			fmt.Println(err)
			continue
		}
		db.Create(&feedMessage)

		fmt.Printf("Message created with id: %d\n", feedMessage.ID)
	}

}

func sendInitMessage(url string, groupId string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{url},
		GroupID: groupId,
		Topic:   topic,
	})
	defer reader.Close()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	msg := kafka.Message{
		Key:   []byte("init-key"),
		Value: []byte("init-message"),
	}

	for {
		err := writer.WriteMessages(context.Background(), msg)
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(10 * time.Second)
	}
}
