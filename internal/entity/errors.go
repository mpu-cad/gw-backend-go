package entity

import "github.com/pkg/errors"

const (
	ErrorParseBody = "can not parse body"
)

var (
	ErrRefreshTokenExpire = errors.New("expire refresh token")
)

type ErrorsRequest struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
