package kitchen

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"pizzeria/internal/kitchen/menu"
	"pizzeria/internal/model"
	"sync"
)

type Menu struct {
	List map[int]menu.Product
}

type Orders struct {
	mu   sync.RWMutex
	list map[uuid.UUID]*model.Order
}

func (o *Orders) Get(id uuid.UUID) (*model.Order, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, ok := o.list[id]

	return order, ok
}

func (o *Orders) Add(order *model.Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.list[order.Id] = order

	_, ok := o.list[order.Id]
	if !ok {
		return fmt.Errorf("order not saved")
	}
	return nil
}

func (o *Orders) Remove(id uuid.UUID) (bool, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.list, id)

	_, ok := o.list[id]
	if !ok {
		return false, fmt.Errorf("order not removed")
	}
	return true, nil
}

type Kitchen struct {
	Menu   Menu
	Orders Orders
	Cooks  []Cook

	InChan  chan *model.Order
	OutChan chan *model.Order
}

func New() *Kitchen {
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
		Orders:  Orders{list: make(map[uuid.UUID]*model.Order), mu: sync.RWMutex{}},
		Menu:    m,
		Cooks:   []Cook{},
		InChan:  make(chan *model.Order, 10),
		OutChan: make(chan *model.Order, 10),
	}

	k.Cooks = append(k.Cooks, Cook{Name: "Андрей", Kitchen: k})

	return k
}

func (k *Kitchen) Work(ctx context.Context) {
	for _, cook := range k.Cooks {
		go cook.Work(ctx)
	}
}

type Cook struct {
	Name    string
	Kitchen *Kitchen
}

func (c *Cook) Work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-c.Kitchen.InChan:
			wg := sync.WaitGroup{}
			for _, product := range order.Products {
				wg.Add(int(product.Quantity))
				for i := int8(1); i <= product.Quantity; i++ {
					fmt.Printf("Повар [\"%s\"] Готовит :%s #%d\n", c.Name, product.Name, i)
					product.Assembling()
					go func() {
						defer wg.Done()
						product.Cooking()
					}()
				}
			}

			wg.Wait()
		}
	}
}

func intPtr(v int) *int {
	return &v
}
