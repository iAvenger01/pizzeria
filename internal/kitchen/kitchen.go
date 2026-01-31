package kitchen

import (
	"context"
	"pizzeria/internal/board"
	"pizzeria/internal/kitchen/menu"
	"pizzeria/internal/model"
	"pizzeria/internal/queue"
	"pizzeria/pkg/logging"
	"time"
)

type Menu struct {
	List map[int]menu.Product `json:"list"`
}

type Kitchen struct {
	db     Storage
	logger *logging.Logger
	Queue  *queue.Queue
	Menu   Menu
	Board  *board.Board
	Cooks  []Cook
	inChan chan *model.Order // Используем только для передачи заказа в горутины, данные "ждут" в очереди, поэтому они небуферизированные
}

func New(logger *logging.Logger, db Storage, board *board.Board) *Kitchen {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(3))
	defer cancel()

	products, _ := db.GetMenu(ctx)
	m := Menu{List: make(map[int]menu.Product)}
	for _, product := range products {
		m.List[product.Id] = product
	}

	k := &Kitchen{
		db:     db,
		logger: logger,
		Board:  board,
		Queue:  queue.NewQueue(50),
		Menu:   m,
		Cooks:  []Cook{},
		inChan: make(chan *model.Order),
	}

	k.Cooks = append(k.Cooks, Cook{Name: "Андрей", kitchen: k})

	return k
}

func (k *Kitchen) Work(ctx context.Context) {
	go func() { // Поставляем горутинам (работникам) заказы в каналы из очереди
		for k.Queue.Len() > 0 || !k.Queue.Closed() {
			order, _ := k.Queue.Dequeue()
			k.inChan <- order
		}
		return
	}()
	for _, cook := range k.Cooks {
		go cook.Work(ctx)
	}
}
