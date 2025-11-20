package orders

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrders(ctx context.Context) ([]Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) Create(ctx context.Context, amount int, status string) (*Order, error) {
	if status == "" {
		status = "new"
	}
	return s.repo.Insert(ctx, amount, status)
}
