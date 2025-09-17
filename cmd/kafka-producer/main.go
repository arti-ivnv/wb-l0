package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"ls-0/arti/order/internal/config"
	"ls-0/arti/order/internal/lib/logger/handlers/slogpretty"
	"math"
	randMath "math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	envLocal = "local"
)

// Mock model for message simulation
type Message struct {
	OrderUuid         string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmId              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"mn_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()    // load config from config file
	log := setupLogger(cfg.Env) // setup logger by env

	log.Info("[Simulator] Starting simulator script...")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "wildberries-topic",
	})

	log.Info("[Simulator] Simulator script started")

	defer writer.Close()

	for {

		log.Info("[Simulator] Seding a new message")

		message := createMessage(log)

		// Send message to kafka
		err := writer.WriteMessages(ctx, kafka.Message{
			Value: []byte(message),
		})

		if err != nil {
			log.Error(fmt.Sprintf("Error while sending message to kafka topic. Error : %s", err.Error()))
		}

		log.Info("[Simulator] Message was sent!")

		time.Sleep(1 * time.Minute)

	}

}

func createMessage(log *slog.Logger) string {

	var msg Message
	var itm Item
	var itm2 Item
	var itm3 Item

	msg.OrderUuid, _ = generateRandomStr(len("b563feb7b2b84b6test"), log)
	msg.TrackNumber, _ = generateRandomStr(len("WBILMTESTTRACK"), log)
	msg.Entry, _ = generateRandomStr(len("WBIL"), log)

	msg.Delivery.Name, _ = generateRandomStr(len("Test Testov"), log)
	msg.Delivery.Phone = fmt.Sprintf("+%s", generateRandomStrNum(len("9720000000"), log))
	msg.Delivery.Zip = generateRandomStrNum(len("2639809"), log)
	msg.Delivery.City, _ = generateRandomStr(len("Kiryat Mozkin"), log)
	msg.Delivery.Address, _ = generateRandomStr(len("Ploshad Mira 15"), log)
	msg.Delivery.Region, _ = generateRandomStr(len("Kraiot"), log)
	msg.Delivery.Email, _ = generateRandomStr(len("test@gmail.com"), log)

	msg.Payment.Transaction, _ = generateRandomStr(len("b563feb7b2b84b6test"), log)
	msg.Payment.RequestId = ""
	msg.Payment.Currency, _ = generateRandomStr(len("USD"), log)
	msg.Payment.Provider, _ = generateRandomStr(len("wbpay"), log)
	msg.Payment.Amount = generateRandomNum(len("1817"))
	msg.Payment.PaymentDT = generateRandomNum(len("11637907727817"))
	msg.Payment.Bank, _ = generateRandomStr(len("alpha"), log)
	msg.Payment.DeliveryCost = generateRandomNum(len("1500"))
	msg.Payment.GoodsTotal = generateRandomNum(len("317"))
	msg.Payment.CustomFee = generateRandomNum(len("10"))

	itm.ChrtId = generateRandomNum(len("9934930"))
	itm.TrackNumber, _ = generateRandomStr(len("WBILMTESTTRACK"), log)
	itm.Price = generateRandomNum(len("453"))
	itm.Rid, _ = generateRandomStr(len("ab4219087a764ae0btest"), log)
	itm.Name, _ = generateRandomStr(len("Mascaras"), log)
	itm.Sale = generateRandomNum(len("2"))
	itm.Size, _ = generateRandomStr(len("0"), log)
	itm.TotalPrice = generateRandomNum(len("11637907727817"))
	itm.NmId = generateRandomNum(len("11637907727817"))
	itm.Brand, _ = generateRandomStr(len("Vivienne"), log)
	itm.Status = generateRandomNum(len("11637907727817"))

	itm2.ChrtId = generateRandomNum(len("9934930"))
	itm2.TrackNumber, _ = generateRandomStr(len("WBILMTESTTRACK"), log)
	itm2.Price = generateRandomNum(len("453"))
	itm2.Rid, _ = generateRandomStr(len("ab4219087a764ae0btest"), log)
	itm2.Name, _ = generateRandomStr(len("Mascaras"), log)
	itm2.Sale = generateRandomNum(len("2"))
	itm2.Size, _ = generateRandomStr(len("0"), log)
	itm2.TotalPrice = generateRandomNum(len("11637907727817"))
	itm2.NmId = generateRandomNum(len("11637907727817"))
	itm2.Brand, _ = generateRandomStr(len("Vivienne"), log)
	itm2.Status = generateRandomNum(len("11637907727817"))

	itm3.ChrtId = generateRandomNum(len("9934930"))
	itm3.TrackNumber, _ = generateRandomStr(len("WBILMTESTTRACK"), log)
	itm3.Price = generateRandomNum(len("453"))
	itm3.Rid, _ = generateRandomStr(len("ab4219087a764ae0btest"), log)
	itm3.Name, _ = generateRandomStr(len("Mascaras"), log)
	itm3.Sale = generateRandomNum(len("2"))
	itm3.Size, _ = generateRandomStr(len("0"), log)
	itm3.TotalPrice = generateRandomNum(len("11637907727817"))
	itm3.NmId = generateRandomNum(len("11637907727817"))
	itm3.Brand, _ = generateRandomStr(len("Vivienne"), log)
	itm3.Status = generateRandomNum(len("111"))

	msg.Items = []Item{itm, itm2, itm3}

	msg.Locale, _ = generateRandomStr(len("en"), log)
	msg.InternalSignature = ""
	msg.CustomerId, _ = generateRandomStr(len("test"), log)
	msg.DeliveryService, _ = generateRandomStr(len("meest"), log)
	msg.Shardkey, _ = generateRandomStr(len("12"), log)
	msg.SmId = generateRandomNum(len("99"))
	msg.DateCreated = "2021-11-26T06:22:19Z" // TODO: TimeStamp
	msg.OofShard, _ = generateRandomStr(len("10"), log)

	jsonMsg, err := json.Marshal(msg)

	if err != nil {
		log.Error("Error marshaling json, err: ", err)
	}

	return string(jsonMsg)
}

func generateRandomStr(length int, log *slog.Logger) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// allocate a byte storage
	byteMsg := make([]byte, length)

	_, err := rand.Read(byteMsg)

	if err != nil {
		log.Error("[Simulator] Error allocating a byte storage")
		return "", err
	}

	for i := 0; i < length; i++ {
		byteMsg[i] = charset[byteMsg[i]%byte(len(charset))]
	}

	return string(byteMsg), nil
}

func generateRandomStrNum(length int, log *slog.Logger) string {
	const charset = "0123456789"

	// allocate a byte storage
	byteMsg := make([]byte, length)

	rand.Read(byteMsg)

	for i := 0; i < length; i++ {
		byteMsg[i] = charset[byteMsg[i]%byte(len(charset))]
	}

	return string(byteMsg)
}

func generateRandomNum(length int) int {
	seed := time.Now().UnixNano()
	source := randMath.NewSource(seed)
	r := randMath.New(source)

	// rand from 1 - 9
	randIntFirstDiv := r.Intn(9) + 1

	maxInt := int(math.Pow10(length+1))*randIntFirstDiv - 1
	minInt := int(math.Pow10(length)) * randIntFirstDiv

	return r.Intn(maxInt-minInt) + minInt
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
