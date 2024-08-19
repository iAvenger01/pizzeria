package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Order struct {
	id      int
	address int
	pizzas  []Pizza
}

type PizzaList struct {
	pizzas []Pizza
}

type Pizza struct {
	name        string
	cookingTime int
}

func main() {
	pizzaList := PizzaList{[]Pizza{
		{"4 сыра", randRange(1, 10)},
		{"Пепперони", randRange(1, 10)},
		{"Дъябло", randRange(1, 10)},
		{"Гавайская", randRange(1, 10)},
		{"Ниндзя", randRange(1, 10)},
		{"Маргарита", randRange(1, 10)},
		{"Пикантная", randRange(1, 10)},
	}}

	var workingHours int
	queueKitchen := make(chan Order, 10)
	queueDelivery := make(chan Order, 10)

	fmt.Println("Сколько часов сегодня работает пиццерия?")
	fmt.Scan(&workingHours)

	// Колл-центр для приема заказов
	go createOrder(queueKitchen, pizzaList, workingHours)

	// Два пицце-мейкера на кухне
	go kitchen(queueKitchen, queueDelivery)
	go kitchen(queueKitchen, queueDelivery)

	// TODO Реализовать возможность добавить несколько курьеров
	delivery(queueDelivery)
}

func createOrder(queueKitchen chan Order, list PizzaList, workingHours int) {
	closedTime := time.Now().Add(time.Hour * time.Duration(workingHours))
	for time.Now().Unix() < closedTime.Unix() {
		countPizzasInOrder := randRange(1, len(list.pizzas))
		pizzas := make([]Pizza, countPizzasInOrder)
		for i := 0; i < countPizzasInOrder; i++ {
			pizzas = append(pizzas, list.pizzas[i])
		}
		order := Order{
			id:      rand.Int(),
			address: randRange(1, 10),
			pizzas:  pizzas,
		}

		queueKitchen <- order
	}

	close(queueKitchen)
}

func kitchen(queueKitchen chan Order, queueDelivery chan Order) {
	for order := range queueKitchen {
		fmt.Printf("Order preparing: %d. Complete: %d / %d \n", order.id, 0, len(order.pizzas))
		for i := 0; i < len(order.pizzas); i++ {
			pizza := order.pizzas[i]
			time.Sleep(time.Second * time.Duration(pizza.cookingTime))
			fmt.Printf("Order is kitchen: %d. Complete: %d / %d \n", order.id, i+1, len(order.pizzas))
		}
		fmt.Println("Order ready: ", order.id)
		queueDelivery <- order
	}
	close(queueDelivery)
}

func delivery(queueDelivery chan Order) {
	for order := range queueDelivery {
		fmt.Println("Order delivering: ", order.id)
		time.Sleep(time.Second * time.Duration(order.address))
		fmt.Println("Order delivered: ", order.id)
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
