package menu

type Item interface {
	Cooking() (bool, error)
}
type Product struct {
	Id             int     `json:"id" db:"id"`
	Key            string  `json:"key" db:"key"`
	Name           string  `json:"name" db:"name"`
	Price          float64 `json:"price" db:"price"`
	CookingTime    int     `json:"-" db:"cooking_time"`
	AssemblingTime *int    `json:"-" db:"assembling_time"`
}
