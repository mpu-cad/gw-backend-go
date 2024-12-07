package configure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mpu-cad/gw-backend-go/internal/configs"
)

func Postgres(ctx context.Context, cfg configs.Postgres) *pgxpool.Pool {
	pool, err := pgxpool.Connect(ctx, cfg.String())

	if err != nil {
		panic("no connect to database")
	}

	return pool
}
