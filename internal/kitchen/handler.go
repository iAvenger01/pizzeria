package kitchen

import (
	"github.com/gofiber/fiber/v2"
	"pizzeria/pkg/logging"
)

const (
	kitchenURL string = "/kitchen"
)

type Handler struct {
	logger  logging.Logger
	kitchen *Kitchen // TODO Заменить на сервис/прослойку для управления кухней не в хандлере
}

func NewHandler(logger logging.Logger, kitchen *Kitchen) Handler {
	return Handler{logger: logger, kitchen: kitchen}
}

func (h *Handler) Register(router fiber.Router) {
	kitchenRouter := router.Group(kitchenURL)
	kitchenRouter.Get("/menu", h.getMenu)
	kitchenRouter.Get("/cooks", h.getCooks)
	kitchenRouter.Get("/orders", h.getOrders)
}

func (h *Handler) getMenu(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	h.logger.Debugln("Get menu from kitchen handler")
	return c.Status(200).JSON(h.kitchen.Menu.List)
}

func (h *Handler) getCooks(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	h.logger.Debugln("Get cooks from kitchen handler")
	return c.Status(200).JSON(h.kitchen.Cooks)
}

func (h *Handler) getOrders(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	h.logger.Debugln("Get orders from kitchen handler")
	return c.Status(200).JSON(h.kitchen.Board.List())
}
