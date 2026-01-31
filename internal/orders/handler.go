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
	logger       *logging.Logger
	orderService Service
}

func NewHandler(logger *logging.Logger, orderService Service) Handler {
	return Handler{logger: logger, orderService: orderService}
}

func (h *Handler) Register(router fiber.Router) {
	router.Get(orderURL, h.getOrder)
	router.Post(ordersURL, h.createOrder)
}

func (h *Handler) getOrder(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	order, err := h.orderService.Get(c.Context(), uuid.MustParse(c.Params("uuid")))
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to get order: %v", err))

		return c.Status(500).SendString("Whoops, something went wrong")
	}
	return c.Status(200).JSON(order)
}

func (h *Handler) createOrder(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	dto := model.OrderDTO{}
	err := c.BodyParser(&dto)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		return c.Status(400).SendString("Request body could not be parsed")
	}
	order, err := h.orderService.Create(c.Context(), dto)
	if err != nil {
		return c.Status(500).SendString("Whoops, something went wrong")
	}
	c.Set(fiber.HeaderLocation, fmt.Sprintf("%s://%sapi/v1/%s/%s", c.Protocol(), c.Hostname(), ordersURL, order.Id.String()))
	return c.Status(201).JSON(order)
}
