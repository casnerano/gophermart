package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/internal/service/account"
	"github.com/casnerano/yandex-gophermart/pkg/logger"
)

type credentialRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Account struct {
	accountService *account.Account
	logger         logger.Logger
}

func NewAccount(service *account.Account, logger logger.Logger) *Account {
	return &Account{accountService: service, logger: logger}
}

func (a *Account) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credential := credentialRequest{}

		err := json.NewDecoder(r.Body).Decode(&credential)
		defer r.Body.Close()

		validPayload := func() bool {
			return credential.Login != "" && credential.Password != ""
		}

		if err != nil || !validPayload() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := a.accountService.SignUp(r.Context(), credential.Login, credential.Password)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExist) {
				w.WriteHeader(http.StatusConflict)
				return
			}

			a.logger.Error("Sing up error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.logger.Info(fmt.Sprintf("Successful sign up user \"%s\"", credential.Login))
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w.WriteHeader(http.StatusOK)
	}
}

func (a *Account) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credential := credentialRequest{}

		err := json.NewDecoder(r.Body).Decode(&credential)
		defer r.Body.Close()

		validPayload := func() bool {
			return credential.Login != "" && credential.Password != ""
		}

		if err != nil || !validPayload() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := a.accountService.SignIn(r.Context(), credential.Login, credential.Password)
		if err != nil {
			if errors.Is(err, account.ErrIncorrectCredentials) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			a.logger.Error("Sign in error.", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.logger.Info(fmt.Sprintf("Successful sign in user \"%s\"", credential.Login))
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w.WriteHeader(http.StatusOK)
	}
}
