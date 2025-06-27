package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"wb-project1/model"

	_ "github.com/lib/pq"
)

func InitPostgres() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=orders_db sslmode=disable"
	return sql.Open("postgres", connStr)
}

func InitSchema(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS orders (
        order_uid TEXT PRIMARY KEY,
        data JSONB NOT NULL
    );`
	_, err := db.Exec(query)
	return err
}

func SaveOrder(db *sql.DB, order model.Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга: %w", err)
	}

	query := `
    INSERT INTO orders (order_uid, data)
    VALUES ($1, $2)
    ON CONFLICT (order_uid) DO UPDATE
    SET data = EXCLUDED.data;`

	_, err = db.Exec(query, order.OrderUID, jsonData)
	return err
}

func GetOrderByID(db *sql.DB, id string) (string, error) {
	var data string
	row := db.QueryRow("SELECT data FROM orders WHERE order_uid = $1", id)
	err := row.Scan(&data)
	if err != nil {
		return "", err
	}
	return data, nil
}

func LoadAllOrders(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT order_uid, data FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[string]string)
	for rows.Next() {
		var id, data string
		if err := rows.Scan(&id, &data); err != nil {
			return nil, err
		}
		orders[id] = data
	}
	return orders, nil
}
