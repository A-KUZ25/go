package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func newDB() (*sql.DB, error) {
	// username:password@protocol(address)/dbname?param=value
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=true&loc=Local"

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
	// db.SetConnMaxLifetime(time.Minute * 5) // время жизни одного коннекта

	return db, nil
}
