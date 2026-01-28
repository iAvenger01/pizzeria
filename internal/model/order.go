package model

import (
	"github.com/google/uuid"
	"pizzeria/internal/kitchen/menu"
)

type OrderDTO struct {
	Address  int16        `json:"address"`
	Products map[int]int8 `json:"products"`
}

type Order struct {
	Id       uuid.UUID       `json:"id" db:"id"`
	Status   string          `json:"status" db:"status"`
	Address  int16           `json:"address" db:"address"`
	Products []*menu.Product `json:"products" db:"-"`
}
