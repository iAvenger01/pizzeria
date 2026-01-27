package menu

import (
	"fmt"
	"time"
)

type Item interface {
	Cooking() (bool, error)
}
type Product struct {
	Status         string  `json:"status" db:"-"`
	Key            string  `json:"key" db:"key"`
	Name           string  `json:"name" db:"name"`
	Price          float64 `json:"price" db:"price"`
	Quantity       int8    `json:"quantity" db:"quantity"`
	CookingTime    int     `json:"-" db:"cooking_time"`
	AssemblingTime *int    `json:"-" db:"assembling_time"`
}

func (p *Product) Assembling() (bool, error) {
	if p.AssemblingTime != nil {
		return false, fmt.Errorf("assembling time not implemented")
	}
	time.Sleep(time.Second * time.Duration(*p.AssemblingTime))
	p.Status = "assembled"
	return true, nil
}

func (p *Product) Cooking() (bool, error) {
	time.Sleep(time.Second * time.Duration(p.CookingTime))
	p.Status = "cooked"
	return true, nil
}
