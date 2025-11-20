package dbConn

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN is not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	// Проверим, что реально можем подключиться
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(50) // максимум открытых коннектов
	db.SetMaxIdleConns(25) // сколько может "болтаться" в простое

	return db, nil
}
