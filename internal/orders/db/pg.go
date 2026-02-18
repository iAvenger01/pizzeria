package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pizzeria/internal/constants"
	"pizzeria/internal/errors"
	"pizzeria/internal/model"
	"pizzeria/internal/orders"
	"pizzeria/pkg/logging"
)

const (
	getOrderQuery = `
		SELECT id, status, address FROM orders WHERE id = $1 LIMIT 1
		`

	getOrderItemsQuery = `
		SELECT p.id, oi.status, p.key, p.name, oi.price, p.assembling_time, p.cooking_time, oi.quantity FROM order_item oi
    	LEFT JOIN products p ON oi.product_id = p.id WHERE oi.order_id = $1 ORDER BY oi.order_id
    `
	insertOrderQuery = `
		INSERT INTO orders (id, address, status) VALUES ($1, $2, $3)
	`
)

var _ orders.Storage = &db{}

type db struct {
	logger *logging.Logger
	pool   *pgxpool.Pool
}

func New(logger *logging.Logger, pool *pgxpool.Pool) orders.Storage {
	return &db{logger: logger, pool: pool}
}

func (d *db) Create(ctx context.Context, orderDTO model.OrderDTO) (uuid.UUID, error) {
	tx, _ := d.pool.Begin(ctx)
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, errors.ErrToCreateUUID
	}
	_, err = tx.Exec(ctx, insertOrderQuery, id, orderDTO.Address, orderDTO.Status)

	if err != nil {
		_ = tx.Rollback(ctx)
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToInsertOrder, err))
		return uuid.UUID{}, fmt.Errorf(constants.ErrToInsertOrder)
	}

	var rows [][]any
	for product, quantity := range orderDTO.Products {
		rows = append(rows, []any{id, product, quantity, 0}) // TODO Придумать элегантный способ передать цену товара
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"order_item"},
		[]string{"order_id", "product_id", "quantity", "price"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToInsertOrderItem, err))
		_ = tx.Rollback(ctx)
		return uuid.UUID{}, err
	}

	_ = tx.Commit(ctx)

	return id, nil
}

func (d *db) FindOne(ctx context.Context, id uuid.UUID) (model.Order, error) {
	tx, _ := d.pool.Begin(ctx)

	rows, err := tx.Query(ctx, getOrderQuery, id)
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToGetOrder, err))
		return model.Order{}, fmt.Errorf(constants.ErrToGetOrder)
	}
	defer rows.Close()
	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToParseOrderInStruct, err))
		return order, fmt.Errorf(constants.ErrToParseOrderInStruct)
	}

	rows, err = tx.Query(ctx, getOrderItemsQuery, id)
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToGetOrderItem, err))
		return order, fmt.Errorf(constants.ErrToGetOrderItem)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.OrderItem])
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToParseOrderItemInStruct, err))
		return order, fmt.Errorf(constants.ErrToParseOrderItemInStruct)
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
