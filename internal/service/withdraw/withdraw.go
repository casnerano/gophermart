package withdraw

import (
	"context"

	"github.com/casnerano/yandex-gophermart/internal/model"
	"github.com/casnerano/yandex-gophermart/internal/repository"
)

type Withdraw struct {
	users     repository.User
	withdraws repository.Withdraw
}

func New(users repository.User, withdraws repository.Withdraw) *Withdraw {
	return &Withdraw{users: users, withdraws: withdraws}
}

func (w *Withdraw) Add(ctx context.Context, orderNumber string, amount float64, userUUID string) (*model.Withdraw, error) {
	return w.withdraws.Add(ctx, orderNumber, amount, userUUID)
}

func (w *Withdraw) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Withdraw, error) {
	return w.withdraws.FindAllByUserUUID(ctx, userUUID)
}

func (w *Withdraw) TotalWithdrawnByUserUUID(ctx context.Context, userUUID string) (float64, error) {
	return w.withdraws.TotalWithdrawnByUserUUID(ctx, userUUID)
}
