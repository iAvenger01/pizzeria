package menu

import "time"

type Pizza struct {
	Product
	Name           string  `json:"name"`
	Cost           float64 `json:"cost"`
	CookingTime    int     `json:"-"`
	AssemblingTime int     `json:"-"`
}

func (p *Pizza) Assembling() bool {
	time.Sleep(time.Second * time.Duration(p.AssemblingTime))
	p.Status = "assembled"
	return true
}

func (p *Pizza) Cooking() (bool, error) {
	time.Sleep(time.Second * time.Duration(p.CookingTime))
	p.Status = "cooked"
	return true, nil
}
