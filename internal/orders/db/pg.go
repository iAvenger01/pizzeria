package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pizzeria/internal/constants"
	"pizzeria/internal/kitchen/menu"
	"pizzeria/internal/model"
	"pizzeria/internal/orders"
	"pizzeria/pkg/logging"
)

const (
	getOrderItemsQuery = `
		SELECT m.key, m.name, m.price, m.assembling_time, m.cooking_time, oi.quantity FROM order_item oi
    	LEFT JOIN menu m ON oi.menu_item_id = m.id WHERE oi.order_id = $1 ORDER BY oi.order_id
    `
	insertOrderQuery = `
		INSERT INTO orders (id, address) VALUES ($1, $2)
	`
)

var _ orders.Storage = &db{}

type db struct {
	logger logging.Logger
	pool   *pgxpool.Pool
}

func New(logger logging.Logger, pool *pgxpool.Pool) orders.Storage {
	return &db{logger: logger, pool: pool}
}

func (d *db) Create(ctx context.Context, orderDTO model.OrderDTO) (uuid.UUID, error) {
	tx, _ := d.pool.Begin(ctx)
	id, err := uuid.NewV7()
	if err != nil {
		d.logger.Error("Error creating uuid: %v", err)
		return uuid.UUID{}, fmt.Errorf("error creating uuid: %v", err)
	}
	_, err = tx.Exec(ctx, insertOrderQuery, id, orderDTO.Address)

	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.UUID{}, err
	}

	var rows [][]any
	for product, quantity := range orderDTO.Products {
		rows = append(rows, []any{id, product, quantity})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"order_item"},
		[]string{"order_id", "menu_item_id", "quantity"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.UUID{}, err
	}

	_ = tx.Commit(ctx)

	return id, nil
}

func (d *db) FindOne(ctx context.Context, id uuid.UUID) (model.Order, error) {
	tx, _ := d.pool.Begin(ctx)

	rows, err := tx.Query(ctx, "SELECT id, address FROM orders WHERE id = $1", id)
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToGetOrder, err))
		return model.Order{}, fmt.Errorf("%s: %v", constants.ErrToGetOrder, err)
	}
	defer rows.Close()
	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		d.logger.Error(err.Error())
		return order, fmt.Errorf("%s: %v", constants.ErrToParseOrderInStruct, err)
	}

	rows, err = tx.Query(ctx, getOrderItemsQuery, id)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[menu.Product])
	if err != nil {
		d.logger.Error(constants.ErrToParseOrderInStruct, err)
		return order, fmt.Errorf("%s: %v", constants.ErrToParseOrderInStruct, err)
	}
	for _, product := range products {
		order.Products = append(order.Products, product)
	}

	return order, nil
}

func (d *db) Update(ctx context.Context, order model.OrderDTO) error {
	//TODO implement me
	panic("implement me")
}
