package kitchen

import (
	"context"
	"fmt"
	"sync"
)

type Cook struct {
	Name    string `json:"name"`
	kitchen *Kitchen
}

func (c *Cook) Work(ctx context.Context) {
	c.kitchen.logger.Info(fmt.Sprintf("[cook][%s] work started", c.Name))
	for {
		select {
		case <-ctx.Done():
			c.kitchen.logger.Info(fmt.Sprintf("[cook][%s] work finished", c.Name))
			return
		case order := <-c.kitchen.inChan:
			c.kitchen.logger.Info(fmt.Sprintf("[order][%s] started cooking", order.Id.String()))
			order.Status = "cooking"
			wg := sync.WaitGroup{}
			for _, product := range order.Products {
				wg.Add(int(product.Quantity))
				for i := int8(1); i <= product.Quantity; i++ {
					c.kitchen.logger.Debug(fmt.Sprintf("[order][%s] assembling: %s #%d", order.Id.String(), product.Key, i))
					product.Assembling()
					c.kitchen.logger.Debug(fmt.Sprintf("[order][%s] assembled: %s #%d", order.Id.String(), product.Key, i))
					go func() {
						defer wg.Done()
						c.kitchen.logger.Debug(fmt.Sprintf("[order][%s] cooking: %s #%d", order.Id.String(), product.Key, i))
						product.Cooking()
						c.kitchen.logger.Debug(fmt.Sprintf("[order][%s] cooked: %s #%d", order.Id.String(), product.Key, i))
					}()
				}
			}

			wg.Wait()
			order.Status = "cooked"
			c.kitchen.logger.Info(fmt.Sprintf("[order][%s] cooked", order.Id.String()))
		}
	}
}
