package model

import (
	"pizzeria/internal/kitchen/menu"
)

type OrderDTO struct {
	Address  int      `json:"address"`
	Products []string `json:"products"`
}

type Order struct {
	Address  int            `json:"address"`
	Products []menu.Product `json:"products"`
}
