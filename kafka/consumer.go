package kafka

import (
	"context"
	"log"
	"wb-project1/config"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func([]byte)

func Consume(topic string, handler MessageHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{config.KafkaBroker},
		Topic:     topic,
		GroupID:   "order-consumer-group",
		Partition: 0,
	})
	if r == nil {
		log.Fatal("Error creating reader")
	}
	defer r.Close()

	log.Println("Kafka consumer listening on topic:", topic)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		handler(msg.Value)
	}
}
