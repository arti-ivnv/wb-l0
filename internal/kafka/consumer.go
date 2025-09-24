package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

func SetUpNewConsumer(ctx context.Context, log *slog.Logger) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "wildberries-topic",
		GroupID: "wb-group",
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error("Error reading a kafka message. err: ", err.Error())
		}
		
		//Save data

		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Error("Error commiting a kafka message. err: ", err.Error())
		}
	}
}
