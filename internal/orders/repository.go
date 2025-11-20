package orders

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]Order, error)
	Insert(ctx context.Context, amount int, status string) (*Order, error)
}
