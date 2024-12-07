package configure

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
)

func Redis(cfg configs.Redis) *redis.Client {
	redisCl := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       int(cfg.DB),
	})

	return redisCl
}
