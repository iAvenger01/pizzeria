package menu

type Product struct {
	Status string
}

type Item interface {
	Cooking() (bool, error)
}
