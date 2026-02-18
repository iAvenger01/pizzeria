package delivery

import (
	"context"
	"github.com/google/uuid"
	"pizzeria/internal/model"
)

type Storage interface {
	Get(ctx context.Context, id uuid.UUID) (model.Courier, error)
	GetAll(ctx context.Context) ([]model.Courier, error)
	Create(ctx context.Context, courier model.CourierDTO) (uuid.UUID, error)
}
