package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/emamulandalib/airbringr-notification/config"
	log "github.com/sirupsen/logrus"
)

func NewSES() (*ses.SES, error) {
	if Sess != nil {
		return ses.New(Sess), nil
	}

	log.Info("New Session should be single.")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Params.AwsRegion),
	})

	if err != nil {
		return nil, err
	}

	Sess = sess
	return ses.New(Sess), nil
}
