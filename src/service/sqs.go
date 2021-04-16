package service

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/emamulandalib/airbringr-notification/config"
	"github.com/emamulandalib/airbringr-notification/dto"
	log "github.com/sirupsen/logrus"
)

var (
	Sess *session.Session
)

type Queue struct {
	Sess session.Session
}

func NewQueue() (*Queue, error) {
	if Sess != nil {
		return &Queue{Sess: *Sess}, nil
	}

	log.Info("New Session")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Params.AwsRegion),
	})

	if err != nil {
		return nil, err
	}

	Sess = sess
	return &Queue{Sess: *sess}, nil
}

func (q *Queue) svc() *sqs.SQS {
	svc := sqs.New(&q.Sess)
	return svc
}

func (q *Queue) List() ([]string, error) {
	urls := []string{}
	svc := q.svc()

	result, err := svc.ListQueues(nil)

	if err != nil {
		return nil, err
	}

	for _, url := range result.QueueUrls {
		urls = append(urls, *url)
	}

	return urls, nil
}

func (q *Queue) GetUrl(name string) (url *string, err error) {
	svc := q.svc()

	res, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &name,
	})

	if err != nil {
		return nil, err
	}

	return res.QueueUrl, err
}

func (q *Queue) Create(name string) error {
	svc := q.svc()

	_, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &name,
		Attributes: map[string]*string{
			"FifoQueue":              aws.String("true"),
			"DelaySeconds":           aws.String("5"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) SendMessage(queueName string, sendSmsDto dto.SendSms) error {
	svc := q.svc()
	url, err := q.GetUrl(queueName)

	if err != nil {
		return err
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         aws.String("sendSms"),
		MessageDeduplicationId: aws.String(q.randString()),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Number": {
				DataType:    aws.String("String"),
				StringValue: aws.String(sendSmsDto.Number),
			},
		},
		MessageBody: aws.String(sendSmsDto.Message),
		QueueUrl:    url,
	})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (q *Queue) ReceiveMessagePeriodic(queueName string) {
	url, err := q.GetUrl(queueName)
	smsService := SmsService{}

	if err != nil {
		panic(err.Error())
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		msgs := q.ReceiveMessage(*url)

		for _, msg := range msgs {
			go smsService.Send(msg)
		}
	}
}

func (q *Queue) ReceiveMessage(url string) []*sqs.Message {
	svc := q.svc()
	res, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &url,
		MaxNumberOfMessages: aws.Int64(10),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
	})

	if err != nil {
		log.Println(err.Error())
	}

	return res.Messages
}

func (q *Queue) DeleteMessage(msg *sqs.Message) {
	svc := q.svc()

	url, err := q.GetUrl(config.Params.SmsQueueName)

	if err != nil {
		log.Println(err.Error())
		return
	}

	svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      url,
		ReceiptHandle: msg.ReceiptHandle,
	})
}
