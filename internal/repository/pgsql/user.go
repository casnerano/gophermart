package pgsql

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casnerano/yandex-gophermart/internal/model"
	"github.com/casnerano/yandex-gophermart/internal/repository"
)

type UserRepository struct {
	pgxpool *pgxpool.Pool
}

func NewUserRepository(pgxpool *pgxpool.Pool) repository.User {
	return &UserRepository{pgxpool}
}

func (p *UserRepository) Add(ctx context.Context, login, password string) (*model.User, error) {
	user := model.User{Login: login, Password: password}
	err := p.pgxpool.QueryRow(
		ctx,
		"insert into users(login, password) values($1, $2) returning uuid, balance, created_at",
		login,
		password,
	).Scan(
		&user.UUID,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = repository.ErrAlreadyExist
		}
		return nil, err
	}

	return &user, nil
}

func (p *UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	user := model.User{Login: login}
	err := p.pgxpool.QueryRow(
		ctx,
		"select uuid, password, balance, created_at from users where login = $1",
		login,
	).Scan(
		&user.UUID,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (p *UserRepository) FindByUUID(ctx context.Context, uuid string) (*model.User, error) {
	user := model.User{UUID: uuid}
	err := p.pgxpool.QueryRow(
		ctx,
		"select login, password, balance, created_at from users where uuid = $1",
		uuid,
	).Scan(
		&user.Login,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
