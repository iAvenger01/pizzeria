package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"pizzeria/internal/kitchen"
	"pizzeria/internal/model"
	"pizzeria/internal/queue"
	"pizzeria/pkg/logging"
)

type service struct {
	logger  *logging.Logger
	kitchen *kitchen.Kitchen
	storage Storage
}

type Service interface {
	Create(ctx context.Context, dto model.OrderDTO) (model.Order, error)
	Get(ctx context.Context, id uuid.UUID) (model.Order, error)
}

var _ Service = &service{}

func NewService(logger *logging.Logger, storage Storage, kitchen *kitchen.Kitchen) (Service, error) {
	return &service{
		storage: storage,
		logger:  logger,
		kitchen: kitchen,
	}, nil
}

func (s *service) Create(ctx context.Context, dto model.OrderDTO) (model.Order, error) {
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

	err = s.kitchen.Queue.Enqueue(order)
	if err != nil {
		if errors.Is(err, queue.ErrQueueClosed) {
			return model.Order{}, fmt.Errorf("kitchen is already closed")
		}
	}
	s.logger.Debug("Sent new order to queue")

	err = s.kitchen.Board.Add(order)
	if err != nil {
		s.logger.Error("Error adding order to kitchen", err)
		return model.Order{}, err
	}
	s.logger.Info("Added new order to kitchen", order.Id.String())

	return *order, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, ok := s.kitchen.Board.Get(id)
	if ok {
		s.logger.Debug("Got order from kitchen", id.String())
		return *order, nil
	}
	return s.storage.FindOne(ctx, id)
}
