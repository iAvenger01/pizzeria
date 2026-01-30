package kitchen

import (
	"context"
	"pizzeria/internal/kitchen/menu"
)

type Storage interface {
	GetMenu(ctx context.Context) ([]menu.Product, error)
}
