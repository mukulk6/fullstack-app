package notifier

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// SendKafkaNotification sends a JSON-encoded message to the Kafka topic
func SendKafkaNotification(topic string, payload interface{}) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"), // Replace with your Kafka broker if different
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	msgBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("❌ Failed to marshal Kafka payload:", err)
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("event-key"),
		Value: msgBytes,
	})
	if err != nil {
		log.Println("❌ Kafka write error:", err)
		return err
	}

	log.Println("✅ Kafka notification sent:", string(msgBytes))
	return nil
}
