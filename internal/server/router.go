package server

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/casnerano/yandex-gophermart/internal/server/handler"
	"github.com/casnerano/yandex-gophermart/internal/server/middleware"
	"github.com/casnerano/yandex-gophermart/internal/service/account"
	"github.com/casnerano/yandex-gophermart/internal/service/balance"
	"github.com/casnerano/yandex-gophermart/internal/service/order"
	"github.com/casnerano/yandex-gophermart/internal/service/withdraw"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

func NewRouter(
	sAccount *account.Account,
	sOrder *order.Order,
	sBalance *balance.Balance,
	sWithdraw *withdraw.Withdraw,
	jwtSecret string,
	logger logger.Logger,
) *chi.Mux {
	accountHandler := handler.NewAccount(sAccount, logger)
	orderHandler := handler.NewOrder(sOrder, logger)
	balanceHandler := handler.NewBalance(sBalance, logger)
	withdrawHandler := handler.NewWithdraw(sWithdraw, logger)

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
		r.Use(middleware.JWTAuthenticator(jwtSecret))
		r.Post("/user/orders", orderHandler.PostUserOrder())
		r.Get("/user/orders", orderHandler.GetUserOrders())
		r.Get("/user/balance", balanceHandler.GetUserSummary())
		r.Post("/user/balance/withdraw", withdrawHandler.PostUserBalanceWithdraw())
		r.Get("/user/withdrawals", withdrawHandler.GetUserWithdrawals())
	})

	router.Mount("/api", router)

	return router
}
