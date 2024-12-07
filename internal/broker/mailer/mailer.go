package mailer

import (
	"github.com/mpu-cad/gw-backend-go/internal/logger"
	"github.com/mpu-cad/gw-backend-go/internal/models"
)

type BrokerMailer struct {
	channel chan models.Gmail
	mailer  mailer
}

func NewBrokerMailer(mailer mailer, channel chan models.Gmail) *BrokerMailer {
	return &BrokerMailer{
		mailer:  mailer,
		channel: channel,
	}
}

func (b *BrokerMailer) Run() {
	for value := range b.channel {
		err := b.mailer.SendEmail(value)
		if err != nil {
			logger.Log.Errorf("err: %v", err)
		}
	}
}
