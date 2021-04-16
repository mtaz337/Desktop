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
	smsQueueName := config.Params.SmsQueueName
	queue, err := service.NewQueue()

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = queue.GetUrl(smsQueueName)

	if err != nil {
		createQueue(queue, smsQueueName)
	}

	go queue.ReceiveMessagePeriodic(smsQueueName)
}

func createQueue(queue *service.Queue, name string) {
	err := queue.Create(name)

	if err != nil {
		log.Fatal(err.Error())
	}
}
