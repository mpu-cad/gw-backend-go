package main

import (
	"context"

	"github.com/mpu-cad/gw-backend-go/internal/app"
	"github.com/mpu-cad/gw-backend-go/internal/configs"
)

func main() {
	ctx := context.Background()
	config := configs.MustConfig(nil)

	application := app.New(config)
	application.Run(ctx)
}
