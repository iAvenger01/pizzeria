package orders

import (
	"pizzeria/internal/kitchen"
	"pizzeria/internal/model"
	"pizzeria/pkg/logging"
)

type service struct {
	logger  logging.Logger
	kitchen kitchen.Kitchen
}

type Service interface {
	CreateOrder(dto model.OrderDTO) (model.Order, error)
}

var _ Service = &service{}

func NewService(logger logging.Logger, kitchen kitchen.Kitchen) (Service, error) {
	return &service{
		logger:  logger,
		kitchen: kitchen,
	}, nil
}

func (s *service) CreateOrder(dto model.OrderDTO) (model.Order, error) {
	order := model.Order{Address: dto.Address}

	for _, p := range dto.Products {
		product := s.kitchen.Menu.List[p]
		order.Products = append(order.Products, product)
	}

	return order, nil
}
