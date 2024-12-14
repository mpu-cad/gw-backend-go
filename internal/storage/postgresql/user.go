package postgresql

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

const (
	SelectUserByID    = `SELECT id, name, surname, last_name, login, email, phone, hash_pass, is_admin, is_blocked, confirm_email FROM "users" WHERE id=$1`
	SelectUserByEmail = `SELECT id, name, surname, last_name, login, email, phone, hash_pass, is_admin, is_blocked, confirm_email FROM "users" WHERE login=$1`
)

type UserRepos struct {
	db *pgxpool.Pool
}

func NewUserRepos(db *pgxpool.Pool) *UserRepos {
	return &UserRepos{
		db: db,
	}
}

func (u *UserRepos) InsertUser(ctx context.Context, user models.User) (*int, error) {
	const (
		query = `
			insert into "users" (name, surname, last_name, login, email, phone, hash_pass) 
			values($1, $2, $3, $4, $5, $6, $7) returning id`
	)

	transaction, err := u.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("can not start transaction, err: %w", err)
	}

	defer func() {
		_ = transaction.Rollback(ctx)
	}()

	var res int
	err = transaction.QueryRow(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.LastName,
		user.Login,
		user.Email,
		user.Phone,
		user.HashPass,
	).Scan(&res)
	if err != nil {
		return nil, fmt.Errorf("can not scan UserRepos for db: %w", err)
	}

	if err = transaction.Commit(ctx); err != nil {
		return nil, fmt.Errorf("can not commit transaction, err: %w", err)
	}

	return &res, nil
}

func (u *UserRepos) SelectUserByID(ctx context.Context, id int) (*models.User, error) {
	return u.getUserFromDB(ctx, SelectUserByID, id)
}

func (u *UserRepos) SelectUserByLogin(ctx context.Context, email string) (*models.User, error) {
	return u.getUserFromDB(ctx, SelectUserByEmail, email)
}

func (u *UserRepos) getUserFromDB(ctx context.Context, query string, arg interface{}) (*models.User, error) {
	var getUser models.User
	err := u.db.QueryRow(ctx, query, arg).
		Scan(
			&getUser.ID,
			&getUser.Name,
			&getUser.Surname,
			&getUser.LastName,
			&getUser.Email,
			&getUser.Login,
			&getUser.Phone,
			&getUser.HashPass,
			&getUser.IsAdmin,
			&getUser.IsBanned,
			&getUser.ConfirmEmail,
		)

	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	return &getUser, nil
}

func (u *UserRepos) ConfirmEmail(ctx context.Context, userID int) error {
	tag, err := u.db.Exec(ctx, `UPDATE users SET confirm_email = true WHERE id = $1`, userID)
	if err != nil {
		return errors.Wrap(err, "update confirm email")
	}

	if !tag.Update() {
		return errors.Wrap(err, "confirm email")
	}

	return nil
}
