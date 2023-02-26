package accrual

import (
	"context"
	"time"

	"github.com/casnerano/yandex-gophermart/internal/service/queue"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type WorkerManager struct {
	rabbitmq *queue.RabbitMQ
	observer *Observer
	logger   logger.Logger
}

func NewWorkerManager(
	rabbitmq *queue.RabbitMQ,
	observer *Observer,
	logger logger.Logger,
) *WorkerManager {
	return &WorkerManager{
		rabbitmq: rabbitmq,
		observer: observer,
		logger:   logger,
	}
}

func (wm *WorkerManager) StartWorker(ctx context.Context) {
	wm.logger.Info("Started accrual queue worker")
	go func() {
		messages, err := wm.rabbitmq.Consume()
		if err != nil {
			wm.logger.Alert("Failed to start accrual worker", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				wm.logger.Info("Stopped accrual worker")
				return
			case message := <-messages:
				ctxObserver, cancelObserver := context.WithTimeout(ctx, time.Minute*15)
				defer cancelObserver()

				wm.logger.Info("Received message from queue and start observe", string(message.Body))
				go wm.observer.Observe(ctxObserver, message)
			}
		}
	}()
}
