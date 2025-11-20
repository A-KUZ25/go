package main

import (
	"log"
	"net/http"
	"test-go/internal/connection"
	"test-go/internal/orders"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/orders", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
	})

	log.Println("server started on :8080")
	http.ListenAndServe(":8080", r)
}
