package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var (
	Params *Config
)

type Config struct {
	Port                  int    `required:"true"`
	AwsRegion             string `required:"true"`
	SmsQueueName          string `required:"true"`
	EmailQueueName        string `required:"true"`
	SmsExternalServiceUrl string `required:"true"`
}

func New() {
	config := Config{}
	err := envconfig.Process("notification", &config)

	if err != nil {
		log.Fatal(err)
	}

	Params = &config
}
