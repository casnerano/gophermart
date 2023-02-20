package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	cfg "github.com/casnerano/yandex-gophermart/internal/config"
	"github.com/casnerano/yandex-gophermart/internal/repository/pgsql"
	srv "github.com/casnerano/yandex-gophermart/internal/server"
	"github.com/casnerano/yandex-gophermart/internal/service/account"
	"github.com/casnerano/yandex-gophermart/internal/service/accrual"
	"github.com/casnerano/yandex-gophermart/internal/service/balance"
	"github.com/casnerano/yandex-gophermart/internal/service/order"
	"github.com/casnerano/yandex-gophermart/internal/service/queue"
	"github.com/casnerano/yandex-gophermart/internal/service/withdraw"
	log "github.com/casnerano/yandex-gophermart/pkg/logger"
)

func main() {

	// Init configuration
	config, err := cfg.New()
	if err != nil {
		panic(err)
	}

	// Init logger
	logger := log.New()
	defer func() {
		if err = logger.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	logLevel := log.LogLevelDebug
	if config.App.ENV == cfg.AppEnvDev {
		logLevel = log.LogLevelWarning
	}

	logger.AddHandler(
		log.NewStdOutHandler(
			log.NewTextFormatter(),
			logLevel,
			true,
		),
	)

	// Database connection
	connection, err := pgxpool.New(context.Background(), config.Database.DSN)
	if err != nil {
		panic(err.Error())
	}

	defer connection.Close()

	// Rabbitmq connection
	rabbitmq, err := queue.NewRabbitMQ(config.Accrual.Queue.DSN, "accrual", "accrual")
	if err != nil {
		logger.Alert("Failed initialization rabbitmq", err)
		os.Exit(1)
	}
	defer rabbitmq.Close()

	// Repositories
	userRepository := pgsql.NewUserRepository(connection)
	orderRepository := pgsql.NewOrderRepository(connection)
	withdrawRepository := pgsql.NewWithdrawRepository(connection)

	// Services
	sAccount := account.New(userRepository, config.App.Secret)
	sOrder := order.New(orderRepository, rabbitmq)
	sBalance := balance.New(userRepository, withdrawRepository)
	sWithdraw := withdraw.New(userRepository, withdrawRepository)

	// Initialization accrual system client
	orderObserver := accrual.NewObserver(
		resty.New(),
		config.Accrual.Service.Address,
		config.Accrual.PoolInterval,
		sOrder,
		logger,
	)

	workerManager := accrual.NewWorkerManager(
		rabbitmq,
		orderObserver,
		logger,
	)

	workerManager.StartWorker(context.Background())

	// Starting server and wait signal for graceful shutdown
	router := srv.NewRouter(
		sAccount,
		sOrder,
		sBalance,
		sWithdraw,
		config.App.Secret,
		logger,
	)

	server := srv.New(config.Server.Address, router, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = server.Run(ctx); err != nil {
		logger.Critical("Failed running server", err)
		os.Exit(1)
	}
}
