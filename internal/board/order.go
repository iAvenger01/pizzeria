package board

import (
	"fmt"
	"github.com/google/uuid"
	"pizzeria/internal/model"
	"sync"
)

type Board struct {
	mu   sync.RWMutex
	list map[uuid.UUID]*model.Order
}

func New() *Board {
	return &Board{mu: sync.RWMutex{}, list: make(map[uuid.UUID]*model.Order)}
}

func (b *Board) List() map[uuid.UUID]*model.Order {
	return b.list
}

func (b *Board) Get(id uuid.UUID) (*model.Order, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	order, ok := b.list[id]

	return order, ok
}

func (b *Board) Add(order *model.Order) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.list[order.Id] = order

	_, ok := b.list[order.Id]
	if !ok {
		return fmt.Errorf("order not saved")
	}
	return nil
}

func (b *Board) Remove(id uuid.UUID) (bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.list, id)

	_, ok := b.list[id]
	if !ok {
		return false, fmt.Errorf("order not removed")
	}
	return true, nil
}
