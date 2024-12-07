package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

const (
	SelectUserByID    = `SELECT * FROM "user" WHERE id=$1`
	SelectUserByEmail = `SELECT * FROM "user" WHERE email=$1`
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
			insert into "user" (name, surname, email, phone, hash_pass) 
			values($1, $2, $3, $4, $5) returning id`
	)

	transaction, err := u.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("can not start transaction, err: %w", err)
	}

	defer func() {
		_ = transaction.Rollback(ctx)
	}()

	var res int
	err = transaction.QueryRow(ctx, query, user.Name, user.Surname, user.Email, user.Phone, user.HashPass).Scan(&res)
	if err != nil {
		return nil, fmt.Errorf("can not scan UserRepos for db: %w", err)
	}

	if err = transaction.Commit(ctx); err != nil {
		return nil, fmt.Errorf("can not commit transaction, err: %w", err)
	}

	return &res, nil
}

func (u *UserRepos) SelectUserByID(ctx context.Context, id int) (*models.User, error) {
	var getUser models.User
	err := u.db.QueryRow(ctx, SelectUserByID, id).
		Scan(
			&getUser.ID,
			&getUser.Name,
			&getUser.Email,
			&getUser.HashPass,
			&getUser.IsAdmin)

	if err != nil {
		return nil, fmt.Errorf("can not scan UserRepos for create db: %w", err)
	}

	return &getUser, nil
}

func (u *UserRepos) SelectUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var getUser models.User
	err := u.db.QueryRow(ctx, SelectUserByEmail, email).Scan(
		&getUser.ID,
		&getUser.Name,
		&getUser.Surname,
		&getUser.Email,
		&getUser.Phone,
		&getUser.HashPass,
		&getUser.IsAdmin,
		&getUser.IsBanned)

	if err != nil {
		return nil, fmt.Errorf("can not select UserRepos by email: %w", err)
	}

	return &getUser, nil
}
