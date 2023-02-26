package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/internal/server/middleware"
	"github.com/casnerano/yandex-gophermart/internal/service/order"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type Order struct {
	orderService *order.Order
	logger       logger.Logger
}

func NewOrder(service *order.Order, logger logger.Logger) *Order {
	return &Order{orderService: service, logger: logger}
}

func (o *Order) PostUserOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, ok := middleware.GetUserUUID(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		orderNumber, err := io.ReadAll(r.Body)
		if err != nil || len(orderNumber) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = o.orderService.Add(r.Context(), string(orderNumber), userUUID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrOrderIncorrectNumber):
				w.WriteHeader(http.StatusUnprocessableEntity)
			case errors.Is(err, order.ErrAlreadyUploaded):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, order.ErrAlreadyUploadedByAnother):
				w.WriteHeader(http.StatusConflict)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			o.logger.Error("Failed to add order", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (o *Order) GetUserOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, ok := middleware.GetUserUUID(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		orders, err := o.orderService.FindAllByUserUUID(r.Context(), userUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			o.logger.Error("Failed find user orders", err)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		bOrders, err := json.Marshal(orders)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			o.logger.Error("Failed marshaller user orders", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(bOrders))
	}
}
