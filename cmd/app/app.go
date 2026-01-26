package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mikhail-bigun/fiberlogrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"pizzeria/internal/config"
	"pizzeria/internal/kitchen"
	"pizzeria/internal/orders"
	"pizzeria/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("Logger initialized")

	cfg, err := config.New()
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("Configuration initialized")

	ctx := context.Background()

	http.Handle("/metrics", promhttp.Handler())

	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      cfg.App.Name,
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
					fiberlogrus.AttachKeyTag(fiberlogrus.TagReqHeader, "Request-Id"),
				},
			},
		),
	)

	api := app.Group("/api/v1")

	k := kitchen.New()
	k.Work(ctx)
	orderService, _ := orders.NewService(logger, k)
	orderHandler := orders.Handler{Logger: logger, OrderService: orderService}
	orderHandler.Register(api)

	app.Get("/menu", func(c *fiber.Ctx) error {
		return c.JSON(k.Menu.List)
	})

	logger.Fatal(app.ListenTLS(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port), "./certs/cert.pem", "./certs/key.pem"))

}
