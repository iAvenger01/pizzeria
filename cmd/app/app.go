package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mikhail-bigun/fiberlogrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"pizzeria/internal/config"
	kitchenPkg "pizzeria/internal/kitchen"
	"pizzeria/internal/orders"
	"pizzeria/internal/orders/db"
	pg "pizzeria/pkg/db"
	"pizzeria/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Debugln("Logger initialized")

	cfg, err := config.New()
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Debugln("Configuration initialized")

	ctx := context.Background()

	pgx, err := pg.New(ctx, cfg)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Debugln("Database initialized")

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

	kitchen := kitchenPkg.New()
	kitchen.Work(ctx)
	orderStorage := db.New(logger, pgx)
	orderService, _ := orders.NewService(logger, orderStorage, kitchen)
	orderHandler := orders.Handler{Logger: logger, OrderService: orderService}
	orderHandler.Register(api)

	app.Get("/menu", func(c *fiber.Ctx) error {
		return c.JSON(kitchen.Menu.List)
	})

	http.Handle("/metrics", promhttp.Handler())

	logger.Fatal(app.ListenTLS(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port), "./certs/cert.pem", "./certs/key.pem"))
}
