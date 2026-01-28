package orders

import (
	"context"
	"github.com/google/uuid"
	"pizzeria/internal/kitchen"
	"pizzeria/internal/model"
	"pizzeria/pkg/logging"
)

type service struct {
	logger  logging.Logger
	kitchen *kitchen.Kitchen
	storage Storage
}

type Service interface {
	CreateOrder(ctx context.Context, dto model.OrderDTO) (model.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error)
}

var _ Service = &service{}

func NewService(logger logging.Logger, storage Storage, kitchen *kitchen.Kitchen) (Service, error) {
	return &service{
		storage: storage,
		logger:  logger,
		kitchen: kitchen,
	}, nil
}

func (s *service) CreateOrder(ctx context.Context, dto model.OrderDTO) (model.Order, error) {
	id, err := s.storage.Create(ctx, dto)
	if err != nil {
		return model.Order{}, err
	}

	order := &model.Order{Id: id, Address: dto.Address}

	for menuItemId, quantity := range dto.Products {
		product := s.kitchen.Menu.List[menuItemId]
		product.Quantity = quantity
		order.Products = append(order.Products, &product)
	}

	err = s.kitchen.Orders.Add(order)
	if err != nil {
		s.logger.Error("Error adding order to kitchen", err)
		return model.Order{}, err
	}
	s.logger.Info("Added new order to kitchen", order.Id.String())

	s.kitchen.InChan <- order
	s.logger.Debug("successfully send order to queue")

	return *order, nil
}

func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, ok := s.kitchen.Orders.Get(id)
	if ok {
		s.logger.Debug("successfully get order from kitchen", id.String())
		return *order, nil
	}
	return s.storage.FindOne(ctx, id)
}
