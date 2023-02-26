package account

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/casnerano/yandex-gophermart/internal/repository"
	"github.com/casnerano/yandex-gophermart/internal/service/token"
)

var (
	ErrIncorrectCredentials = errors.New("incorrect credentials")
)

type Account struct {
	users  repository.User
	secret string
}

func New(users repository.User, secret string) *Account {
	return &Account{users: users, secret: secret}
}

func (a *Account) SignIn(ctx context.Context, login, password string) (string, error) {
	user, err := a.users.FindByLogin(ctx, login)
	if err != nil {
		return "", ErrIncorrectCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrIncorrectCredentials
	}

	jwtToken, err := token.NewJWT(user.UUID, a.secret)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (a *Account) SignUp(ctx context.Context, login, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := a.users.Add(ctx, login, string(hashedPassword))
	if err != nil {
		return "", err
	}

	jwtToken, err := token.NewJWT(user.UUID, a.secret)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
