package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mikhail-bigun/fiberlogrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"pizzeria/internal/kitchen"
	"pizzeria/internal/model"
	"pizzeria/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("Logger initialized")

	ctx := context.Background()

	k := kitchen.New()
	k.Work(ctx)

	http.Handle("/metrics", promhttp.Handler())

	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
	})
	app.Use(
		fiberlogrus.New(
			fiberlogrus.Config{
				Logger: logger.Logger,
				Tags: []string{
					// add method field
					fiberlogrus.TagLatency,
					fiberlogrus.TagMethod,
					// add status field
					fiberlogrus.TagStatus,
					fiberlogrus.TagRoute,
					fiberlogrus.TagPath,
					// add value from locals
					fiberlogrus.AttachKeyTag(fiberlogrus.TagReqHeader, "requestid"),
				},
			},
		),
	)

	app.Get("/menu", func(c *fiber.Ctx) error {
		return c.JSON(k.Menu.List)
	})

	app.Post("/orders", func(c *fiber.Ctx) error {
		order := model.Order{}
		err := c.BodyParser(&order)
		if err != nil {
			logger.Error(err)
		}
		k.InChan <- order

		return c.Status(fiber.StatusCreated).JSON(order)
	})

	app.ListenTLS("pizzeria.local:8080", "./certs/cert.pem", "./certs/key.pem")

}
