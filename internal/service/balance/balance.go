package balance

import (
	"context"

	"github.com/casnerano/yandex-gophermart/internal/repository"
)

type Balance struct {
	users     repository.User
	withdraws repository.Withdraw
}

type Summary struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func New(users repository.User, withdraws repository.Withdraw) *Balance {
	return &Balance{users: users, withdraws: withdraws}
}

func (b *Balance) GetSummaryByUserUUID(ctx context.Context, userUUID string) (*Summary, error) {
	user, err := b.users.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	withdrawn, err := b.withdraws.TotalWithdrawnByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	summary := Summary{
		Current:   user.Balance,
		Withdrawn: withdrawn,
	}
	return &summary, nil
}
