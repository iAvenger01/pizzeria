package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	boardPkg "pizzeria/internal/board"
	"pizzeria/internal/model"
	"pizzeria/internal/queue"
	"pizzeria/pkg/logging"
	"time"
)

type Delivery struct {
	db       Storage
	logger   *logging.Logger
	Board    *boardPkg.Board
	Queue    *queue.Queue
	inChan   chan *model.Order
	Couriers map[uuid.UUID]*Courier
}

func New(logger *logging.Logger, db Storage, board *boardPkg.Board) *Delivery {
	_, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	defer cancel()

	d := &Delivery{
		logger:   logger,
		Board:    board,
		Queue:    queue.NewQueue(50),
		inChan:   make(chan *model.Order),
		Couriers: map[uuid.UUID]*Courier{},
	}

	couriers, err := db.GetAll(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf("couriers not available: %w", err)
			logger.Fatal(err.Error())
		}
	}

	for _, courier := range couriers {
		d.Couriers[courier.Id] = &Courier{
			Id:       courier.Id,
			Name:     courier.Name,
			WorkTime: courier.WorkTime,
			Bag: &Bag{
				places: make([]*model.Order, 0, courier.BagSize),
				size:   courier.BagSize,
			},
			deliveryService: d,
		}
	}

	return d
}

func (d *Delivery) Work(ctx context.Context) {
	go func() { // Поставляем горутинам (работникам) заказы в каналы из очереди
		for d.Queue.Len() > 0 || !d.Queue.Closed() {
			order, _ := d.Queue.Dequeue()
			d.inChan <- order
		}
		return
	}()
	for _, courier := range d.Couriers {
		go courier.Work(ctx)
	}
}
