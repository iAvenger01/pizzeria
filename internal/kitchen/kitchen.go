package kitchen

import (
	"context"
	"pizzeria/internal/kitchen/menu"
	"pizzeria/internal/model"
	"sync"
)

type Menu struct {
	List map[string]menu.Item `json:"list"`
}

type Kitchen struct {
	Menu    Menu
	Cooks   []Cook
	InChan  chan model.Order
	OutChan chan model.Order
}

func New() *Kitchen {
	list := map[string]menu.Pizza{
		"four_cheese": {Name: "4 сыра", Cost: 390.0, CookingTime: 2, AssemblingTime: 15},
		"pepperoni":   {Name: "Пепперони", Cost: 589.0, AssemblingTime: 2, CookingTime: 15},
		"diablo":      {Name: "Дъябло", Cost: 539.0, AssemblingTime: 5, CookingTime: 15},
		"hawaiian":    {Name: "Гавайская", Cost: 480.0, AssemblingTime: 4, CookingTime: 15},
		"ninja":       {Name: "Ниндзя", Cost: 580.0, AssemblingTime: 4, CookingTime: 15},
		"margarita":   {Name: "Маргарита", Cost: 350.0, AssemblingTime: 3, CookingTime: 15},
		"spicy":       {Name: "Пикантная", Cost: 650.0, AssemblingTime: 5, CookingTime: 15},
		"bbq":         {Name: "Барбекю", Cost: 690.0, AssemblingTime: 4, CookingTime: 15},
	}

	var m map[string]menu.Item

	for key, item := range list {
		m[key] = item
	}

	k := &Kitchen{
		Menu:    m,
		Cooks:   []Cook{},
		InChan:  make(chan model.Order, 10),
		OutChan: make(chan model.Order, 10),
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
			wg.Add(len(order.Products))
			for _, _ = range order.Products {
				//originalPizza := c.Kitchen.Menu.List[productKey]
				//fmt.Printf("Повар [\"%s\"] Готовит :%s\n", c.Name, originalPizza.Name)
				//originalPizza.Assembling()
				go func() {
					defer wg.Done()
					//originalPizza.Cook()
				}()
			}

			wg.Wait()
		}
	}
}
