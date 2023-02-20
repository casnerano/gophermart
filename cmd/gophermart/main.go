package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	cfg "github.com/casnerano/yandex-gophermart/internal/config"
	"github.com/casnerano/yandex-gophermart/internal/repository/pgsql"
	srv "github.com/casnerano/yandex-gophermart/internal/server"
	"github.com/casnerano/yandex-gophermart/internal/server/handler"
	"github.com/casnerano/yandex-gophermart/internal/server/middleware"
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

	// Account dependencies
	userRepository := pgsql.NewUserRepository(connection)
	accountService := account.New(userRepository, config.App.Secret)
	accountHandler := handler.NewAccount(accountService, logger)

	// Order dependencies
	orderRepository := pgsql.NewOrderRepository(connection)
	orderService := order.New(orderRepository, rabbitmq)
	orderHandler := handler.NewOrder(orderService, logger)

	// Balance dependencies
	withdrawRepository := pgsql.NewWithdrawRepository(connection)
	balanceService := balance.New(userRepository, withdrawRepository)
	balanceHandler := handler.NewBalance(balanceService, logger)

	// Withdraw dependencies
	withdrawService := withdraw.New(userRepository, withdrawRepository)
	withdrawHandler := handler.NewWithdraw(withdrawService, logger)

	// Router
	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.JSONContentType)

	// Public routes
	router.Group(func(r chi.Router) {
		r.Post("/user/register", accountHandler.SignUp())
		r.Post("/user/login", accountHandler.SignIn())
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator(config.App.Secret))
		r.Post("/user/orders", orderHandler.PostUserOrder())
		r.Get("/user/orders", orderHandler.GetUserOrders())
		r.Get("/user/balance", balanceHandler.GetUserSummary())
		r.Post("/user/balance/withdraw", withdrawHandler.PostUserBalanceWithdraw())
		r.Get("/user/withdrawals", withdrawHandler.GetUserWithdrawals())
	})

	router.Mount("/api", router)

	// Initialization accrual system client
	orderObserver := accrual.NewObserver(
		resty.New(),
		config.Accrual.Service.Address,
		config.Accrual.PoolInterval,
		orderService,
		logger,
	)

	workerManager := accrual.NewWorkerManager(
		rabbitmq,
		orderObserver,
		logger,
	)

	workerManager.StartWorker(context.Background())

	// Starting server and wait signal for graceful shutdown
	server := srv.New(config.Server.Address, router, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = server.Run(ctx); err != nil {
		logger.Critical("Failed running server", err)
		os.Exit(1)
	}
}
