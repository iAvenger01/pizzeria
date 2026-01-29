package delivery

import (
	"context"
	"fmt"
	boardPkg "pizzeria/internal/board"
	"pizzeria/internal/model"
	"pizzeria/internal/queue"
)

type Delivery struct {
	Board    *boardPkg.Board
	Queue    *queue.Queue
	inChan   chan *model.Order
	Couriers []Courier
}

func New(board *boardPkg.Board) *Delivery {
	d := &Delivery{
		Board:    board,
		Queue:    queue.NewQueue(50),
		inChan:   make(chan *model.Order),
		Couriers: []Courier{},
	}

	d.Couriers = append(d.Couriers, Courier{Name: "Андрей", delivery: d})

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

type Courier struct {
	Name     string `json:"name"`
	delivery *Delivery
}

func (c *Courier) Work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-c.delivery.inChan:
			order.Status = "delivering"
			fmt.Printf("Курьер [%s] Доставляет: заказ %s\n", c.Name, order.Id.String())
			order.Status = "delivered"
		}
	}
}
