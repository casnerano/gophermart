package order

import (
	"context"
	"errors"

	"github.com/casnerano/yandex-gophermart/internal/model"
	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/internal/service/queue"
)

var (
	ErrAlreadyUploaded          = errors.New("already uploaded")
	ErrAlreadyUploadedByAnother = errors.New("already uploaded by another")
)

type Order struct {
	orders   repository.Order
	rabbitmq *queue.RabbitMQ
}

func New(orders repository.Order, rabbitmq *queue.RabbitMQ) *Order {
	return &Order{orders: orders, rabbitmq: rabbitmq}
}

func (o *Order) Add(ctx context.Context, number, userUUID string) (*model.Order, error) {
	order, err := o.orders.Add(ctx, number, userUUID)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExist) {
			foundOrder, foundErr := o.orders.FindByNumber(ctx, number)
			if foundErr != nil {
				return nil, foundErr
			}
			if foundOrder.UserUUID == userUUID {
				return nil, ErrAlreadyUploaded
			}
			return nil, ErrAlreadyUploadedByAnother
		}
		return nil, err
	}

	_ = o.rabbitmq.Publish(ctx, []byte(order.Number))
	return order, nil
}

func (o *Order) FindByNumber(ctx context.Context, number string) (*model.Order, error) {
	return o.orders.FindByNumber(ctx, number)
}

func (o *Order) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Order, error) {
	return o.orders.FindAllByUserUUID(ctx, userUUID)
}

func (o *Order) AccrueByNumber(ctx context.Context, number string, status model.OrderStatus, accrual float64) (*model.Order, error) {
	return o.orders.AccrueByNumber(ctx, number, status, accrual)
}
