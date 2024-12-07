package logger

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"

	"github.com/mpu-cad/gw-backend-go/internal/configs"
)

var Log = logrus.New()

func InitLogger(cfg configs.Logger) {
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "02-01-2006 15:04:05.000",
		PrettyPrint:     true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
		},
	})

	logrus.SetOutput(colorable.NewColorableStdout())

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	err := Log.Level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		Log.Panicf("failed to set log level: %v", err)
	}

	Log.Infof("log level set to %v", cfg.LogLevel)
}
