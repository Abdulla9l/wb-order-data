package repository

import (
	"database/sql"
	"fmt"
	"log"
	"wb-order-data/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveOrder(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin error: %v", err)
	}
	defer tx.Rollback()

	orderQuery := `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature,
	                  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	               ON CONFLICT (order_uid) DO UPDATE SET 
	                  track_number = EXCLUDED.track_number,
	                  entry = EXCLUDED.entry,
	                  locale = EXCLUDED.locale,
	                  internal_signature = EXCLUDED.internal_signature,
	                  customer_id = EXCLUDED.customer_id,
	                  delivery_service = EXCLUDED.delivery_service,
	                  shardkey = EXCLUDED.shardkey,
	                  sm_id = EXCLUDED.sm_id,
	                  date_created = EXCLUDED.date_created,
	                  oof_shard = EXCLUDED.oof_shard`

	_, err = tx.Exec(orderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID,
		order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("save order error: %v", err)
	}

	deliveryQuery := `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
	                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	                  ON CONFLICT (order_uid) DO UPDATE SET
	                     name = EXCLUDED.name,
	                     phone = EXCLUDED.phone,
	                     zip = EXCLUDED.zip,
	                     city = EXCLUDED.city,
	                     address = EXCLUDED.address,
	                     region = EXCLUDED.region,
	                     email = EXCLUDED.email`

	_, err = tx.Exec(deliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("get delivery error: %v", err)
	}

	paymentQuery := `INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, 
	                     payment_dt, bank, delivery_cost, goods_total, custom_fee)
	                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	                 ON CONFLICT (order_uid) DO UPDATE SET
	                    transaction = EXCLUDED.transaction,
	                    request_id = EXCLUDED.request_id,
	                    currency = EXCLUDED.currency,
	                    provider = EXCLUDED.provider,
	                    amount = EXCLUDED.amount,
	                    payment_dt = EXCLUDED.payment_dt,
	                    bank = EXCLUDED.bank,
	                    delivery_cost = EXCLUDED.delivery_cost,
	                    goods_total = EXCLUDED.goods_total,
	                    custom_fee = EXCLUDED.custom_fee`

	_, err = tx.Exec(paymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("get payment error: %v", err)
	}

	_, err = tx.Exec("DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return fmt.Errorf("delete items error: %v", err)
	}

	itemQuery := `INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
	                   sale, size, total_price, nm_id, brand, status)
	              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		_, err = tx.Exec(itemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)

		if err != nil {
			return fmt.Errorf("get items error: %v", err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	orderQuery := `SELECT order_uid, track_number, entry, locale, internal_signature,
	                      customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard 
	               FROM orders WHERE order_uid = $1`

	err := r.db.QueryRow(orderQuery, id).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID,
		&order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get order error: %v", err)
	}

	var delivery models.Delivery
	deliveryQuery := `SELECT name, phone, zip, city, address, region, email 
	                  FROM delivery WHERE order_uid = $1`

	err = r.db.QueryRow(deliveryQuery, id).Scan(
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("get delivery error: %v", err)
	}
	order.Delivery = delivery

	var payment models.Payment
	paymentQuery := `SELECT transaction, request_id, currency, provider, amount, payment_dt, 
	                        bank, delivery_cost, goods_total, custom_fee 
	                 FROM payment WHERE order_uid = $1`

	err = r.db.QueryRow(paymentQuery, id).Scan(
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
		&payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("get payment error: %v", err)
	}
	order.Payment = payment

	itemsQuery := `SELECT chrt_id, track_number, price, rid, name, sale, size, 
	                      total_price, nm_id, brand, status 
	               FROM items WHERE order_uid = $1`

	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("get items error: %v", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item error: %v", err)
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	query := `SELECT order_uid FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get orders error: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, fmt.Errorf("scan order_uid error: %v", err)
		}

		order, err := r.GetOrderByID(orderUID)
		if err != nil {
			log.Printf("Warning: failed to get order %s: %v", orderUID, err)
			continue
		}
		if order != nil {
			orders = append(orders, *order)
		}
	}

	return orders, nil
}
