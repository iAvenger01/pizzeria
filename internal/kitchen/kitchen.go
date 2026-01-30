package kitchen

import (
	"context"
	"pizzeria/internal/board"
	"pizzeria/internal/kitchen/menu"
	"pizzeria/internal/model"
	"pizzeria/internal/queue"
	"pizzeria/pkg/logging"
)

type Menu struct {
	List map[int]menu.Product `json:"list"`
}

type Kitchen struct {
	logger *logging.Logger
	Queue  *queue.Queue
	Menu   Menu
	Board  *board.Board
	Cooks  []Cook
	inChan chan *model.Order // Используем только для передачи заказа в горутины, данные "ждут" в очереди, поэтому они небуферизированные
}

func New(logger *logging.Logger, board *board.Board) *Kitchen {
	m := Menu{List: map[int]menu.Product{
		1: {Name: "4 сыра", Price: 390.0, AssemblingTime: intPtr(2), CookingTime: 15},
		2: {Name: "Пепперони", Price: 589.0, AssemblingTime: intPtr(2), CookingTime: 15},
		3: {Name: "Дъябло", Price: 539.0, AssemblingTime: intPtr(5), CookingTime: 15},
		4: {Name: "Гавайская", Price: 480.0, AssemblingTime: intPtr(4), CookingTime: 15},
		5: {Name: "Ниндзя", Price: 580.0, AssemblingTime: intPtr(4), CookingTime: 15},
		6: {Name: "Маргарита", Price: 350.0, AssemblingTime: intPtr(3), CookingTime: 15},
		7: {Name: "Пикантная", Price: 650.0, AssemblingTime: intPtr(5), CookingTime: 15},
		8: {Name: "Барбекю", Price: 690.0, AssemblingTime: intPtr(4), CookingTime: 15},
	}}

	k := &Kitchen{
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

func intPtr(v int) *int {
	return &v
}
