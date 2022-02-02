package app

import (
	"github.com/emamulandalib/airbringr-notification/config"
	"github.com/emamulandalib/airbringr-notification/service"
	log "github.com/sirupsen/logrus"
)

func (app *App) Bootstrap() {
	app.createQueueOrFail()
}

func (app *App) createQueueOrFail() {
	smsSvc := service.SmsService{}
	emailSvc := service.EmailService{}

	smsQueueName := config.Params.SmsQueueName
	emailQueueName := config.Params.EmailQueueName

	createQueue(smsQueueName)
	createQueue(emailQueueName)

	go smsSvc.ReceiveMessagePeriodic(smsQueueName)
	go emailSvc.ReceiveMessagePeriodic(emailQueueName)
}

func createQueue(name string) {
	queue, err := service.NewQueue()

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = queue.GetUrl(name)

	if err == nil {
		return
	}

	err = queue.Create(name)

	if err != nil {
		log.Fatal(err.Error())
	}
}
