package orders

import (
	"context"
	"pizzeria/internal/model"
)

type Storage interface {
	Create(ctx context.Context, order model.Order) (string, error)
	FindOne(ctx context.Context, uuid string) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
	Delete(ctx context.Context, uuid string) error
}
