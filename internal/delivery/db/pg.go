package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pizzeria/internal/delivery"
	"pizzeria/internal/errors"
	"pizzeria/internal/model"
	"pizzeria/pkg/logging"
)

const (
	insertEmployeeQuery = `
		INSERT INTO employees (id, first_name, last_name, status, work_time, employment_date)
		VALUES (@id, @first_name, @last_name, @status, @work_time, NOW());
	`
	insertCourierQuery = `
		INSERT INTO couriers (employee_id, bag_size) VALUES (@employee_id, @bag_size)
	`
	getOneCourierQuery = `
		SELECT e.id, CONCAT(e.last_name, ' ', e.first_name) as name, e.work_time, c.bag_size
		FROM couriers c LEFT JOIN employees e on e.id = c.employee_id WHERE c.employee_id = $1 LIMIT 1
	`
	getAllCouriersQuery = `
		SELECT e.id, CONCAT(e.last_name, ' ', e.first_name) as name, e.work_time, c.bag_size
		FROM couriers c LEFT JOIN employees e on e.id = c.employee_id
	`
)

var _ delivery.Storage = &db{}

type db struct {
	logger *logging.Logger
	pool   *pgxpool.Pool
}

func New(logger *logging.Logger, pool *pgxpool.Pool) delivery.Storage {
	return &db{logger: logger, pool: pool}
}

func (d *db) Get(ctx context.Context, id uuid.UUID) (model.Courier, error) {
	rows, err := d.pool.Query(ctx, getOneCourierQuery, id)
	if err != nil {
		return model.Courier{}, fmt.Errorf("failed to get courier: %w", err)
	}
	defer rows.Close()

	courier, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Courier])
	if err != nil {
		return courier, fmt.Errorf("failed to parse courier: %w", err)
	}

	return courier, nil
}

func (d *db) GetAll(ctx context.Context) ([]model.Courier, error) {
	rows, err := d.pool.Query(ctx, getAllCouriersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get couriers: %w", err)
	}
	defer rows.Close()

	couriers, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Courier])
	if err != nil {
		return couriers, fmt.Errorf("failed to parse couriers in struct: %w", err)
	}

	return couriers, nil
}

func (d *db) Create(ctx context.Context, dto model.CourierDTO) (uuid.UUID, error) {
	tx, _ := d.pool.Begin(ctx)

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, errors.ErrToCreateUUID
	}

	args := pgx.NamedArgs{
		"id":         id.String(),
		"first_name": dto.FirstName,
		"last_name":  dto.LastName,
		"status":     "active",
		"work_time":  dto.WorkTime,
	}
	_, err = tx.Exec(ctx, insertEmployeeQuery, args)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert employee: %w", err)
	}

	_, err = tx.Exec(ctx, insertCourierQuery, pgx.NamedArgs{"employee_id": id.String(), "bag_size": dto.BagSize})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert courier: %w", err)
	}

	return uuid.UUID{}, nil
}
