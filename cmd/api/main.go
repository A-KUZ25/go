package main

import (
	"log"
	"net/http"
	"test-go/internal/connection"
	"test-go/internal/orders"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	db, err := connection.NewDB()
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	repo, err := orders.NewMySQLRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	service := orders.NewService(repo)
	handler := orders.NewHandler(service)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.List(w, r)
		case http.MethodPost:
			handler.Create(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
