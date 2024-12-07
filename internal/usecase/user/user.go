package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type UCUser struct {
	user userRepos
}

func NewUCUser(user userRepos) *UCUser {
	return &UCUser{
		user: user,
	}
}

func (u *UCUser) Registration(ctx context.Context, request models.User) (*int, error) {
	pass, err := u.createHashPassword(request.HashPass)
	if err != nil {
		return nil, err
	}

	request.HashPass = pass

	id, err := u.user.InsertUser(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("can not insert user, err: %w", err)
	}

	return id, nil
}

func (u *UCUser) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.user.SelectUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = u.compareHashAndPassword(user.HashPass, password)
	if err != nil {
		return nil, errors.New("wrong password or email")
	}

	return user, nil
}

func (u *UCUser) createHashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("create hashed password was failed: %v", err.Error())
	}

	return string(hashPassword), nil
}

func (u *UCUser) compareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("compare is wrong, err: %w", err)
	}

	return nil
}
