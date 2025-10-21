package database

import (
    "database/sql"
    "fmt"
    "log"
    "wb-order-data/internal/config"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
    )

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to open DB connection: %v", err)
    }

    // Проверка подключения
    if err = DB.Ping(); err != nil {
        log.Fatalf("Failed to ping DB: %v", err)
    }

    log.Println("Connected to PostgreSQL!")
}
