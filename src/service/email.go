package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/emamulandalib/airbringr-notification/config"
	"github.com/emamulandalib/airbringr-notification/dto"
	log "github.com/sirupsen/logrus"
	"time"
)

type EmailService struct{}

func (svc *EmailService) EnQueue(sendEmail *dto.SendEmail) error {
	message := sendEmail.Message
	queue, err := NewQueue()

	if err != nil {
		log.Error(err.Error())
		return err
	}

	data, _ := json.Marshal(sendEmail.Data)
	msgAttrs := map[string]*sqs.MessageAttributeValue{
		"From": {
			DataType:    aws.String("String"),
			StringValue: aws.String(sendEmail.From),
		},
		"To": {
			DataType:    aws.String("String"),
			StringValue: aws.String(sendEmail.To),
		},
		"Subject": {
			DataType:    aws.String("String"),
			StringValue: aws.String(sendEmail.Subject),
		},
		"TemplateCode": {
			DataType:    aws.String("String"),
			StringValue: aws.String(sendEmail.TemplateCode),
		},
		"Data": {
			DataType:    aws.String("String"),
			StringValue: aws.String(string(data)),
		},
	}

	if sendEmail.CC != nil {
		msgAttrs["CC"] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(*sendEmail.CC),
		}
	}

	if sendEmail.BCC != nil {
		msgAttrs["BCC"] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(*sendEmail.BCC),
		}
	}

	err = queue.SendMessage(config.Params.EmailQueueName, "sendEmail", message, msgAttrs)

	if err != nil {
		return err
	}

	return nil
}

func (svc *EmailService) Send(msg *sqs.Message) {
	from := *msg.MessageAttributes["From"].StringValue
	to := *msg.MessageAttributes["To"].StringValue
	subject := *msg.MessageAttributes["Subject"].StringValue
	tmplCode := *msg.MessageAttributes["TemplateCode"].StringValue
	data := *msg.MessageAttributes["Data"].StringValue
	charset := "UTF-8"

	var tmplData interface{}
	_ = json.Unmarshal([]byte(data), &tmplData)
	html, err := GenerateTpl(tmplCode, tmplData)

	if err != nil {
		log.Error(err.Error())
		return
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	sesSvc, err := NewSES()

	if err != nil {
		log.Error(err.Error())
		return
	}

	_, err = sesSvc.SendEmail(input)

	if err != nil {
		log.Error(err.Error())
		return
	}

	q, _ := NewQueue()
	q.DeleteMessage(config.Params.EmailQueueName, msg)
}

func (svc *EmailService) ReceiveMessagePeriodic(queueName string) {
	q, err := NewQueue()

	if err != nil {
		log.Panic(err.Error())
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		msgs := q.ReceiveMessage(queueName)

		for _, msg := range msgs {
			go svc.Send(msg)
		}
	}
}
