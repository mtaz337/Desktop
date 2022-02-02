package service

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/emamulandalib/airbringr-notification/config"
	log "github.com/sirupsen/logrus"
)

type Queue struct {
	Sess session.Session
}

func NewQueue() (*Queue, error) {
	if Sess != nil {
		return &Queue{Sess: *Sess}, nil
	}

	log.Info("New Session should be single.")

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

func (q *Queue) SendMessage(name string, groupID string, msgBody string, msgAttrs map[string]*sqs.MessageAttributeValue) error {
	svc := q.svc()
	url, err := q.GetUrl(name)

	if err != nil {
		return err
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         aws.String(groupID),
		MessageDeduplicationId: aws.String(q.randString()),
		MessageAttributes:      msgAttrs,
		MessageBody:            aws.String(msgBody),
		QueueUrl:               url,
	})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (q *Queue) ReceiveMessagePeriodic(name string) {
	url, err := q.GetUrl(name)
	smsService := SmsService{}

	if err != nil {
		log.Panic(err.Error())
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

func (q *Queue) ReceiveMessage(name string) []*sqs.Message {
	url, err := q.GetUrl(name)

	if err != nil {
		log.Panic(err.Error())
	}

	svc := q.svc()
	res, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            url,
		MaxNumberOfMessages: aws.Int64(10),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
	})

	if err != nil {
		log.Error(err.Error())
	}

	return res.Messages
}

func (q *Queue) DeleteMessage(name string, msg *sqs.Message) {
	svc := q.svc()

	url, err := q.GetUrl(name)

	if err != nil {
		log.Error(err.Error())
		return
	}

	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      url,
		ReceiptHandle: msg.ReceiptHandle,
	})

	if err != nil {
		log.Error(err.Error())
	}
}
