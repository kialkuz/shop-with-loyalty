package db

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/infrastructure/storage"
	"kialkuz/shop-with-loyalty/internal/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *storage.DB
}

func NewUserRepository(db *storage.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CheckExistLogin(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&exists,
		)
	}, "SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&user.ID,
			&user.Login,
			&user.Password,
		)
	}, "SELECT id, login, password FROM users WHERE login = $1", login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByToken(ctx context.Context, jti string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&user.ID,
			&user.Login,
		)
	}, `SELECT u.id, u.login
		FROM users u
		JOIN active_tokens a_t ON a_t.user_id = u.id
		WHERE a_t.jti = $1`,
		jti)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkgErrors.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) AddNewUser(
	ctx context.Context,
	user model.User,
) error {
	return r.db.Exec(
		ctx,
		"INSERT INTO users (id, login, password) VALUES ($1, $2, $3)",
		user.ID,
		user.Login,
		user.Password,
	)
}

func (r *UserRepository) AddNewUserTx(
	ctx context.Context,
	tx pgx.Tx,
	user model.User,
) error {
	_, err := tx.Exec(
		ctx,
		"INSERT INTO users (id, login, password) VALUES ($1, $2, $3)",
		user.ID,
		user.Login,
		user.Password,
	)

	return err
}
