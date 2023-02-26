package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/casnerano/yandex-gophermart/internal/server/middleware"
	"github.com/casnerano/yandex-gophermart/internal/service/balance"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type Balance struct {
	balanceService *balance.Balance
	logger         logger.Logger
}

func NewBalance(service *balance.Balance, logger logger.Logger) *Balance {
	return &Balance{balanceService: service, logger: logger}
}

func (b *Balance) GetUserSummary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, ok := middleware.GetUserUUID(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		summary, err := b.balanceService.GetSummaryByUserUUID(r.Context(), userUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			b.logger.Error("Failed get user balance summary", err)
			return
		}

		bSummary, err := json.Marshal(summary)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			b.logger.Error("Failed marshaller user balance summary", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(bSummary))
	}
}
