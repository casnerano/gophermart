package pgsql

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casnerano/yandex-gophermart/internal/model"
	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/pkg/luhn"
)

type OrderRepository struct {
	pgxpool *pgxpool.Pool
}

func NewOrderRepository(pgxpool *pgxpool.Pool) repository.Order {
	return &OrderRepository{pgxpool}
}

func (p *OrderRepository) Add(ctx context.Context, number, userUUID string) (*model.Order, error) {
	if !luhn.Checksum(number) {
		return nil, repository.ErrOrderIncorrectNumber
	}

	order := model.Order{Number: number, UserUUID: userUUID}
	err := p.pgxpool.QueryRow(
		ctx,
		"insert into orders(number, status, user_uuid) values($1, $2, $3) returning uuid, accrual, uploaded_at",
		number,
		model.OrderStatusNew,
		userUUID,
	).Scan(
		&order.UUID,
		&order.Accrual,
		&order.UploadedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = repository.ErrAlreadyExist
		}
		return nil, err
	}

	return &order, nil
}

func (p *OrderRepository) FindByNumber(ctx context.Context, number string) (*model.Order, error) {
	order := model.Order{Number: number}
	err := p.pgxpool.QueryRow(
		ctx,
		"select uuid, status, accrual, user_uuid, uploaded_at from orders where number = $1",
		number,
	).Scan(
		&order.UUID,
		&order.Status,
		&order.Accrual,
		&order.UserUUID,
		&order.UploadedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (p *OrderRepository) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	rows, err := p.pgxpool.Query(
		ctx,
		"select uuid, number, status, accrual, user_uuid, uploaded_at from orders where user_uuid = $1 order by uploaded_at",
		userUUID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		order := &model.Order{}
		err = rows.Scan(
			&order.UUID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UserUUID,
			&order.UploadedAt,
		)
		if err == nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (p *OrderRepository) AccrueByNumber(ctx context.Context, number string, status model.OrderStatus, accrual float64) (*model.Order, error) {
	order := model.Order{
		Number:  number,
		Accrual: accrual,
		Status:  status,
	}

	tx, err := p.pgxpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		"update orders set status = $1, accrual = $2 where number = $3 returning uuid, user_uuid, uploaded_at",
		status,
		accrual,
		number,
	).Scan(
		&order.UUID,
		&order.UserUUID,
		&order.UploadedAt,
	)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		ctx,
		"update users set balance = balance + $1 where uuid = $2",
		accrual,
		order.UserUUID,
	)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
