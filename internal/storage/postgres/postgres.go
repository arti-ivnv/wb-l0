package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"ls-0/arti/order/internal/config"
	"ls-0/arti/order/internal/storage"
	"sync"

	"github.com/jackc/pgx/v5"

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
		return fmt.Errorf("error unmarshaling json data: %w", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(
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

	_, err = tx.Exec(
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

	_, err = tx.Exec(
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
		_, err = tx.Exec(
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
			return fmt.Errorf("failed to save item: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (ps *PostgresStorage) GetAll(ctx context.Context) []storage.Order {

	conn, err := ps.pool.Acquire(ctx)
	if err != nil {
		fmt.Println("failed to acquire connection: ", err)
		return nil
	}

	defer conn.Release()

	var ordersToSave []storage.Order

	orderRows, _ := conn.Query(ctx, "SELECT * FROM orders")

	for orderRows.Next() {
		var ord storage.Order

		if err := orderRows.Scan(
			&ord.OrderUuid,
			&ord.TrackNumber,
			&ord.Entry,
			&ord.Locale,
			&ord.InternalSignature,
			&ord.CustomerId,
			&ord.DeliveryService,
			&ord.Shardkey,
			&ord.SmId,
			&ord.DateCreated,
			&ord.OofShard,
		); err != nil {
			return nil
		}
		ordersToSave = append(ordersToSave, ord)
	}

	orderRows.Close()

	for i := range ordersToSave {
		deliveryRow, _ := conn.Query(ctx, "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1", ordersToSave[i].OrderUuid)
		delivery, err := pgx.CollectExactlyOneRow(deliveryRow, pgx.RowToStructByName[storage.Delivery])
		if err != nil {
			fmt.Println("delivery_rec_err: ", err)
		}
		fmt.Println(delivery)

		ordersToSave[i].Delivery = delivery
		deliveryRow.Close()

		paymentRow, _ := conn.Query(ctx, "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1", ordersToSave[i].OrderUuid)
		payment, err := pgx.CollectExactlyOneRow(paymentRow, pgx.RowToStructByName[storage.Payment])
		if err != nil {
			fmt.Println("payment_rec_err: ", err)
		}
		ordersToSave[i].Payment = payment
		paymentRow.Close()
	}

	for i := range ordersToSave {

		itemRows, err := conn.Query(ctx, "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1", ordersToSave[i].OrderUuid)
		if err != nil {
			fmt.Println(err)
		}
		var orderItems []storage.Item
		for itemRows.Next() {
			var itm storage.Item
			if err := itemRows.Scan(
				&itm.ChrtId,
				&itm.TrackNumber,
				&itm.Price,
				&itm.Rid,
				&itm.Name,
				&itm.Sale,
				&itm.Size,
				&itm.TotalPrice,
				&itm.NmId,
				&itm.Brand,
				&itm.Status,
			); err != nil {
				return nil
			}

			orderItems = append(orderItems, itm)
		}
		ordersToSave[i].Items = orderItems
		itemRows.Close()
	}

	fmt.Println(len(ordersToSave))
	fmt.Println(ordersToSave)

	return ordersToSave
}
