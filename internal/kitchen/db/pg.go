package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pizzeria/internal/constants"
	"pizzeria/internal/kitchen"
	"pizzeria/internal/kitchen/menu"
	"pizzeria/pkg/logging"
)

const (
	getMenuQuery = `
		SELECT id, key, name, price, assembling_time, cooking_time FROM menu
	`
)

var _ kitchen.Storage = &db{}

type db struct {
	logger *logging.Logger
	pool   *pgxpool.Pool
}

func New(logger *logging.Logger, pool *pgxpool.Pool) kitchen.Storage {
	return &db{logger: logger, pool: pool}
}

func (d *db) GetMenu(ctx context.Context) ([]menu.Product, error) {
	rows, err := d.pool.Query(ctx, getMenuQuery)
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToGetMenu, err))
		return nil, err
	}
	defer rows.Close()

	menuItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[menu.Product])
	if err != nil {
		d.logger.Error(fmt.Sprintf("%s: %v", constants.ErrToParseMenuItemsInStruct, err))
		return menuItems, fmt.Errorf(constants.ErrToParseMenuItemsInStruct)
	}

	return menuItems, nil
}
