package orders

import "time"

type Order struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
