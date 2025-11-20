package orders

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetOrders(r.Context())
	if err != nil {
		http.Error(w, "failed: "+err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int    `json:"amount"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	order, err := h.service.Create(r.Context(), input.Amount, input.Status)
	if err != nil {
		http.Error(w, "failed: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
