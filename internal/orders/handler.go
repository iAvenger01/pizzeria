package orders

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"pizzeria/internal/model"
	"pizzeria/pkg/logging"
)

const (
	ordersURL string = "/orders"
	orderURL  string = "/orders/:uuid"
)

type Handler struct {
	Logger       logging.Logger
	OrderService Service
}

func (h *Handler) Register(router fiber.Router) {
	router.Get(orderURL, h.getOrder)
	router.Post(ordersURL, h.createOrder)
}

func (h *Handler) getOrder(c *fiber.Ctx) error {
	order, err := h.OrderService.GetOrder(c.Context(), uuid.MustParse(c.Params("uuid")))
	if err != nil {
		h.Logger.Error("Failed to get order", err)
	}
	return c.Status(200).JSON(order)
}

func (h *Handler) createOrder(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	dto := model.OrderDTO{}
	err := c.BodyParser(&dto)
	if err != nil {
		h.Logger.Error(err)
	}
	order, err := h.OrderService.CreateOrder(c.Context(), dto)
	if err != nil {
		h.Logger.Error(err)
		return c.Status(500).SendString("Whoops, something went wrong")
	}
	c.Set(fiber.HeaderLocation, fmt.Sprintf("%s://%sapi/v1/%s/%s", c.Protocol(), c.Hostname(), ordersURL, order.Id.String()))
	return c.Status(201).JSON(order)
}
