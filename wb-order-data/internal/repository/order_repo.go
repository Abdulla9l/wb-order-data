package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"wb-order-data/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

// Конструктор репозитория
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// SaveOrder сохраняет заказ в БД
func (r *OrderRepository) SaveOrder(order *models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	query := `INSERT INTO orders (order_uid, data) VALUES ($1, $2)
			  ON CONFLICT (order_uid) DO UPDATE SET data = EXCLUDED.data`
	_, err = r.db.Exec(query, order.OrderUID, orderJSON)
	if err != nil {
		return fmt.Errorf("db save error: %v", err)
	}

	return nil
}

// GetOrderByID возвращает заказ по его ID
func (r *OrderRepository) GetOrderByID(id string) (*models.Order, error) {
	query := `SELECT data FROM orders WHERE order_uid = $1`
	row := r.db.QueryRow(query, id)

	var orderJSON []byte
	if err := row.Scan(&orderJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // заказ не найден
		}
		return nil, fmt.Errorf("db query error: %v", err)
	}

	var order models.Order
	if err := json.Unmarshal(orderJSON, &order); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}

	return &order, nil
}

// GetAllOrders загружает все заказы (например, для восстановления кеша)
func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	query := `SELECT data FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query error: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var orderJSON []byte
		if err := rows.Scan(&orderJSON); err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		var order models.Order
		if err := json.Unmarshal(orderJSON, &order); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %v", err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}
