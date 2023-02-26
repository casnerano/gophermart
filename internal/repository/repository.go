package repository

import (
	"context"
	"errors"

	"github.com/casnerano/yandex-gophermart/internal/model"
)

var (
	ErrAlreadyExist = errors.New("entity already exists")
	ErrNotFound     = errors.New("entity not found")

	ErrOrderIncorrectNumber     = errors.New("incorrect order number")
	ErrWithdrawNotEnoughBalance = errors.New("not enough balance")
)

type User interface {
	Add(ctx context.Context, login, password string) (*model.User, error)
	FindByLogin(ctx context.Context, login string) (*model.User, error)
	FindByUUID(ctx context.Context, uuid string) (*model.User, error)
}

type Order interface {
	Add(ctx context.Context, number, userUUID string) (*model.Order, error)
	FindByNumber(ctx context.Context, number string) (*model.Order, error)
	FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Order, error)
	AccrueByNumber(ctx context.Context, number string, status model.OrderStatus, accrual float64) (*model.Order, error)
}

type Withdraw interface {
	Add(ctx context.Context, orderNumber string, amount float64, userUUID string) (*model.Withdraw, error)
	FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Withdraw, error)
	TotalWithdrawnByUserUUID(ctx context.Context, userUUID string) (float64, error)
}
