package orders

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pizzeria/internal/model"
	"pizzeria/pkg/logging"
)

const (
	ordersURL string = "/orders/"
	orderURL  string = "/orders/:uuid"
)

type Handler struct {
	Logger       logging.Logger
	OrderService Service
}

func (h *Handler) Register(router fiber.Router) {
	router.Get(orderURL, h.Order)
	router.Post(ordersURL, h.CreateOrder)
}

func (h *Handler) Order(c *fiber.Ctx) error {
	return c.Status(200).SendString("Order is processing")
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Set(fiber.HeaderLocation, fmt.Sprintf("%s/%s", ordersURL, "5555"))
	dto := model.OrderDTO{}
	err := c.BodyParser(&dto)
	if err != nil {
		h.Logger.Error(err)
	}
	order, _ := h.OrderService.CreateOrder(dto)

	return c.Status(201).JSON(order)
}
