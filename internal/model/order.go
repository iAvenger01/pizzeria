package model

import (
	"fmt"
	"github.com/google/uuid"
	"pizzeria/internal/kitchen/menu"
	"time"
)

type OrderDTO struct {
	Address  int16        `json:"address"`
	Products map[int]int8 `json:"products"`
	Status   string       `json:"-"`
}

type Order struct {
	Id       uuid.UUID    `json:"id" db:"id"`
	Status   string       `json:"status" db:"status"`
	Address  int16        `json:"address" db:"address"`
	Products []*OrderItem `json:"products" db:"-"`
}

var _ menu.Item = &OrderItem{}

type OrderItem struct {
	Status   string `json:"status" db:"status"`
	Quantity int8   `json:"quantity" db:"quantity"`
	menu.Product
}

func (i *OrderItem) Assembling() (bool, error) {
	if i.AssemblingTime == nil {
		return false, fmt.Errorf("assembling time not implemented")
	}
	time.Sleep(time.Second * time.Duration(*i.AssemblingTime))
	i.Status = "assembled"
	return true, nil
}

func (i *OrderItem) Cooking() (bool, error) {
	time.Sleep(time.Second * time.Duration(i.CookingTime))
	i.Status = "cooked"
	return true, nil
}
