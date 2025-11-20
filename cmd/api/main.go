package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"test-go/internal/dbConn"
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	db *sql.DB
)

func fetchOrdersFromDB(db *sql.DB) ([]Order, error) {
	rows, err := db.Query(`SELECT id, amount, status, created_at, updated_at FROM orders ORDER BY id DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Order

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Amount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOrders(w, r)
	case http.MethodPost:
		createOrder(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := fetchOrdersFromDB(db)
	if err != nil {
		http.Error(w, "failed to load orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func insertOrder(db *sql.DB, amount int, status string) (*Order, error) {
	res, err := db.Exec(`INSERT INTO orders (amount, status) VALUES (?, ?)`, amount, status)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var o Order
	row := db.QueryRow(`SELECT id, amount, status, created_at, updated_at FROM orders WHERE id = ?`, id)
	if err := row.Scan(&o.ID, &o.Amount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}

	return &o, nil
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int    `json:"amount"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if input.Status == "" {
		input.Status = "new"
	}

	order, err := insertOrder(db, input.Amount, input.Status)
	if err != nil {
		http.Error(w, "failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	var err error
	db, err = dbConn.NewDB()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/orders", ordersHandler)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
