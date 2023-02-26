package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/internal/server/middleware"
	"github.com/casnerano/yandex-gophermart/internal/service/withdraw"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type Withdraw struct {
	withdrawService *withdraw.Withdraw
	logger          logger.Logger
}

func NewWithdraw(service *withdraw.Withdraw, logger logger.Logger) *Withdraw {
	return &Withdraw{withdrawService: service, logger: logger}
}

func (wd *Withdraw) PostUserBalanceWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, ok := middleware.GetUserUUID(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		withdrawRequest := struct {
			Order string  `json:"order"`
			Sum   float64 `json:"sum"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&withdrawRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			wd.logger.Error("Failed withdrawn request unmarshall", err)
			return
		}

		_, err = wd.withdrawService.Add(
			r.Context(),
			withdrawRequest.Order,
			withdrawRequest.Sum,
			userUUID,
		)

		if err != nil {
			switch {
			case errors.Is(err, repository.ErrWithdrawNotEnoughBalance):
				w.WriteHeader(http.StatusPaymentRequired)
			case errors.Is(err, repository.ErrOrderIncorrectNumber):
				w.WriteHeader(http.StatusUnprocessableEntity)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			wd.logger.Error("Failed to add withdraw", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (wd *Withdraw) GetUserWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, ok := middleware.GetUserUUID(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		withdrawals, err := wd.withdrawService.FindAllByUserUUID(r.Context(), userUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			wd.logger.Error("Failed find user withdrawals", err)
			return
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		bWithdrawals, err := json.Marshal(withdrawals)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			wd.logger.Error("Failed marshaller user withdrawals", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(bWithdrawals))
	}
}
