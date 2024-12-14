package redis

import (
	"context"
	"strconv"
	"time"
)

const (
	ttlCode = 15 * time.Minute
)

func (t *TokenRepos) SaveUsersRegistrationCode(ctx context.Context, code string, userID int) {
	t.cli.Set(ctx, "RegistrationCode"+strconv.Itoa(userID), code, ttlCode)
}

func (t *TokenRepos) GetUsersRegistrationCode(ctx context.Context, userID int) string {
	return t.cli.Get(ctx, "RegistrationCode"+strconv.Itoa(userID)).String()
}
