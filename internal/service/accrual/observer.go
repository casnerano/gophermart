package accrual

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/casnerano/yandex-gophermart/internal/model"
	"github.com/casnerano/yandex-gophermart/internal/service/order"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type Data struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type Observer struct {
	client       *resty.Client
	accrualURL   string
	poolInterval int
	orderService *order.Order
	logger       logger.Logger
}

func NewObserver(
	client *resty.Client,
	accrualURL string,
	poolInterval int,
	orderService *order.Order,
	logger logger.Logger,
) *Observer {
	return &Observer{
		client:       client,
		accrualURL:   accrualURL,
		poolInterval: poolInterval,
		orderService: orderService,
		logger:       logger,
	}
}

func (o *Observer) Observe(ctx context.Context, delivery amqp.Delivery) error {
	for {
		data := Data{}
		response, err := o.client.SetBaseURL(o.accrualURL + "/api").
			R().SetResult(&data).Get("/orders/" + string(delivery.Body))
		if err != nil {
			o.logger.Error(fmt.Sprintf("Failed order `\"%s\" observe in accrual system", delivery.Body), err)
			_ = delivery.Nack(false, true)
			return err
		}

		switch response.StatusCode() {
		case http.StatusOK:
			err = o.UpdateOrder(ctx, data)
			if err != nil {
				o.logger.Error(fmt.Sprintf("Failed processing accrual for order `\"%s\" in accrual system", delivery.Body), err)
				_ = delivery.Nack(false, true)
			}
			_ = delivery.Ack(false)
			o.logger.Info(fmt.Sprintf("Successfully processing accrua for order `\"%s\"", delivery.Body))
			return nil
		case http.StatusNoContent:
			o.logger.Error(fmt.Sprintf("Not registered order `\"%s\" in accrual system", delivery.Body), err)
			_ = delivery.Ack(false)
			return nil
		case http.StatusTooManyRequests:
			o.logger.Error(fmt.Sprintf("Too Many Requests for order `\"%s\" in accrual system", delivery.Body), err)

			retryAfter := response.Header().Get("Retry-After")
			if retryAfter == "" {
				_ = delivery.Nack(false, true)
			}

			var retryAfterValue int
			retryAfterValue, err = strconv.Atoi(retryAfter)
			if err != nil {
				retryAfterValue = 1
			}

			select {
			case <-ctx.Done():
				_ = delivery.Nack(false, true)
			case <-time.After(time.Second * time.Duration(retryAfterValue)):
			}
		default:
			o.logger.Error(fmt.Sprintf("Unknown response status for order `\"%s\" in accrual system", delivery.Body), response.Error())
			return errors.New("unknown response status")
		}
	}
}

func (o *Observer) UpdateOrder(ctx context.Context, data Data) error {
	var status model.OrderStatus
	switch data.Status {
	case "INVALID":
		status = model.OrderStatusInvalid
	case "PROCESSING":
		status = model.OrderStatusProcessing
	case "PROCESSED":
		status = model.OrderStatusProcessed
	default:
		return errors.New("unknown order status")
	}

	_, err := o.orderService.AccrueByNumber(ctx, data.Order, status, data.Accrual)
	return err
}
