package delivery

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"pizzeria/internal/model"
	"time"
)

type Courier struct {
	Id              uuid.UUID    `json:"id" db:"id"`
	Name            string       `json:"name" db:"name"`
	WorkTime        int          `json:"work_time" db:"work_time"`
	Bag             *Bag         `json:"bag"`
	CurrentOrder    *model.Order `json:"current_order" db:"current_order"`
	deliveryService *Delivery
}

type Bag struct {
	places []*model.Order // TODO переделать на кол-во продуктов для возможности перевозки нескольких заказов
	size   int
}

func (c *Courier) Work(ctx context.Context) {
	c.deliveryService.logger.Info(fmt.Sprintf("[courier][%s] work started", c.Name))
	for {
		select {
		case <-ctx.Done():
			c.deliveryService.logger.Info(fmt.Sprintf("[courier][%s] work finished: %v", c.Name, ctx.Err()))
			return
		case order := <-c.deliveryService.inChan:
			c.addOrder(order)
			c.deliveryService.logger.Debug(fmt.Sprintf("[courier][%s] added order in bag", c.Name))

			if !c.isFullBag() {
			loop:
				for {
					select {
					case <-time.After(time.Second * 5):
						c.deliveryService.logger.Debug(fmt.Sprintf("[courier][%s] wait more orders", c.Name))
						break loop
					case moreOrder := <-c.deliveryService.inChan:
						c.addOrder(moreOrder)
						c.deliveryService.logger.Debug(fmt.Sprintf("[courier][%s] added order in bag", c.Name))
						if !c.isFullBag() {
							c.deliveryService.logger.Debug(fmt.Sprintf("[courier][%s] bag is full", c.Name))
							break loop
						}
					}
				}
			}
			c.deliveryService.logger.Debug(fmt.Sprintf("[courier][%s] delivery started", c.Name))
			c.delivery()
			c.deliveryService.logger.Debug(fmt.Sprintf("[courier][%s] delivery completed", c.Name))
		}
	}
}

func (c *Courier) addOrder(order *model.Order) {
	c.Bag.places = append(c.Bag.places, order)
}

func (c *Courier) isFullBag() bool {
	return len(c.Bag.places) == c.Bag.size
}

func (c *Courier) delivery() {
	for _, order := range c.Bag.places {
		order.Status = "delivering"
		c.deliveryService.logger.Info(fmt.Sprintf("[order][%s] delivering", order.Id.String()))
		time.Sleep(time.Second * time.Duration(order.Address))
		order.Status = "delivered"
		c.deliveryService.logger.Info(fmt.Sprintf("[order][%s] delivered", order.Id.String()))
	}

	c.Bag.places = c.Bag.places[:0]
}
