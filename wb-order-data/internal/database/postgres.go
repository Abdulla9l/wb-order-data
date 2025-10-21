package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"wb-order-data/internal/config"
)

func Connect(cfg *config.Config) *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("❌ Ошибка открытия подключения: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ PostgreSQL недоступен: %v", err)
	}

	log.Println("✅ Подключено к PostgreSQL!")
	return db
}
