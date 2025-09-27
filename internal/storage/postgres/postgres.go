package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"ls-0/arti/order/internal/config"
	"ls-0/arti/order/internal/storage"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	mu   sync.RWMutex
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) *PostgresStorage {

	dbUrl := cfg.Pg.Url

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil
	}

	// Customize pool config
	dbConfig.MaxConns = cfg.Pg.Pool.MaxConns
	dbConfig.MinConns = cfg.Pg.Pool.MinConns
	dbConfig.MaxConnLifetime = cfg.Pg.Pool.MaxConnLifetime
	dbConfig.HealthCheckPeriod = cfg.Pg.Pool.HealthcheckPeriod
	dbConfig.MinIdleConns = cfg.Pg.Pool.MinIdle

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil
	}

	return &PostgresStorage{pool: pool}
}

func (ps *PostgresStorage) AddOrder(inOrder string, ctx context.Context) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var orderToSave storage.Order

	conn, err := ps.pool.Acquire(ctx)
	if err != nil {
		fmt.Println("failed to acquire connection: ", err)
		return fmt.Errorf("failed to acquire connection: %w", err)
	}

	defer conn.Release()

	jsonOrder := []byte(inOrder)

	err = json.Unmarshal(jsonOrder, &orderToSave)
	if err != nil {
		fmt.Println("Error unmarshaling json data: ", err)
		return fmt.Errorf("Error unmarshaling json data: %w", err)
	}

	_, err = conn.Exec(
		ctx,
		"INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		orderToSave.OrderUuid,
		orderToSave.TrackNumber,
		orderToSave.Entry,
		orderToSave.Locale,
		orderToSave.InternalSignature,
		orderToSave.CustomerId,
		orderToSave.DeliveryService,
		orderToSave.Shardkey,
		orderToSave.SmId,
		orderToSave.DateCreated,
		orderToSave.OofShard,
	)

	if err != nil {
		fmt.Println("failed to save order: ", err)
		return fmt.Errorf("failed to save order: %w", err)
	}

	_, err = conn.Exec(
		ctx,
		"INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES($1, $2, $3, $4, $5, $6, $7, $8)",
		orderToSave.OrderUuid,
		orderToSave.Delivery.Name,
		orderToSave.Delivery.Phone,
		orderToSave.Delivery.Zip,
		orderToSave.Delivery.City,
		orderToSave.Delivery.Address,
		orderToSave.Delivery.Region,
		orderToSave.Delivery.Email,
	)

	if err != nil {
		fmt.Println("failed to save delivery: ", err)
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	_, err = conn.Exec(
		ctx,
		"INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		orderToSave.OrderUuid,
		orderToSave.Payment.Transaction,
		orderToSave.Payment.RequestId,
		orderToSave.Payment.Currency,
		orderToSave.Payment.Provider,
		orderToSave.Payment.Amount,
		orderToSave.Payment.PaymentDT,
		orderToSave.Payment.Bank,
		orderToSave.Payment.DeliveryCost,
		orderToSave.Payment.GoodsTotal,
		orderToSave.Payment.CustomFee,
	)

	if err != nil {
		fmt.Println("failed to save payment: ", err)
		return fmt.Errorf("failed to save payment: %w", err)
	}

	for _, item := range orderToSave.Items {
		_, err = conn.Exec(
			ctx,
			"INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			orderToSave.OrderUuid,
			item.ChrtId,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status,
		)

		if err != nil {
			fmt.Errorf("failed to save item: ", err)
			return fmt.Errorf("failed to save item: %w", err)
		}
	}

	return nil
}
