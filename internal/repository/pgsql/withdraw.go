package pgsql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casnerano/yandex-gophermart/internal/model"
	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/pkg/luhn"
)

type WithdrawRepository struct {
	pgxpool *pgxpool.Pool
}

func NewWithdrawRepository(pgxpool *pgxpool.Pool) repository.Withdraw {
	return &WithdrawRepository{pgxpool}
}

func (w *WithdrawRepository) Add(ctx context.Context, orderNumber string, amount float64, userUUID string) (*model.Withdraw, error) {
	if !luhn.Checksum(orderNumber) {
		return nil, repository.ErrOrderIncorrectNumber
	}

	order := model.Withdraw{OrderNumber: orderNumber, Amount: amount, UserUUID: userUUID}

	tx, err := w.pgxpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var balance float64
	err = tx.QueryRow(
		ctx,
		"select balance from users where uuid = $1",
		userUUID,
	).Scan(&balance)

	if err != nil {
		return nil, err
	}

	if balance < amount {
		return nil, repository.ErrWithdrawNotEnoughBalance
	}

	_, err = tx.Exec(
		ctx,
		"update users set balance = balance - $1 where uuid = $2",
		amount,
		userUUID,
	)

	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		ctx,
		"insert into withdraws(order_number, amount, user_uuid) values($1, $2, $3) returning uuid, processed_at",
		orderNumber,
		amount,
		userUUID,
	).Scan(
		&order.UUID,
		&order.ProcessedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = repository.ErrAlreadyExist
		}
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (w *WithdrawRepository) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Withdraw, error) {
	withdraws := make([]*model.Withdraw, 0)
	rows, err := w.pgxpool.Query(
		ctx,
		"select uuid, order_number, amount, user_uuid, processed_at from withdraws where user_uuid = $1 order by processed_at",
		userUUID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		withdraw := &model.Withdraw{}
		err = rows.Scan(
			&withdraw.UUID,
			&withdraw.OrderNumber,
			&withdraw.Amount,
			&withdraw.UserUUID,
			&withdraw.ProcessedAt,
		)
		if err == nil {
			withdraws = append(withdraws, withdraw)
		}
	}

	return withdraws, nil
}

func (w *WithdrawRepository) TotalWithdrawnByUserUUID(ctx context.Context, userUUID string) (float64, error) {
	var total sql.NullFloat64
	err := w.pgxpool.QueryRow(
		ctx,
		"select sum(amount) from withdraws where user_uuid = $1",
		userUUID,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	if !total.Valid {
		return 0, nil
	}

	return total.Float64, nil
}
