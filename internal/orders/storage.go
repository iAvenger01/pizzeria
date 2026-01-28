package orders

import (
	"context"
	"github.com/google/uuid"
	"pizzeria/internal/model"
)

type Storage interface {
	Create(ctx context.Context, order model.OrderDTO) (uuid.UUID, error)
	Update(ctx context.Context, order model.OrderDTO) error
	FindOne(ctx context.Context, uuid uuid.UUID) (model.Order, error)
}
