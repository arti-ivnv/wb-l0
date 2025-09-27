package kafka

import (
	"context"
	"log/slog"
	"ls-0/arti/order/internal/storage/postgres"
	"ls-0/arti/order/internal/storage/safer"

	"github.com/segmentio/kafka-go"
)

func SetUpNewConsumer(ctx context.Context, storage *postgres.PostgresStorage, sfm *safer.SafeMap, log *slog.Logger) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "wildberries-topic",
		GroupID: "wb-group",
	})

	defer reader.Close()

	for {

		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error("[CONSUMER] Error reading a kafka message. err: ", err.Error())
		} else {
			log.Info("[CONSUMER] Getting a new kafka msg")
		}

		// Save date to db
		storage.AddOrder(string(msg.Value), ctx)

		// Save data to map
		sfm.Put(string(msg.Value), log)

		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Error("[CONSUMER] Error commiting a kafka message. err: ", err.Error())
		} else {
			log.Info("[CONSUMER] Kafka msg was commited")
		}

	}
}
