package service

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/emamulandalib/airbringr-notification/config"
	"github.com/emamulandalib/airbringr-notification/dto"
	log "github.com/sirupsen/logrus"
)

type SmsService struct{}

func (svc *SmsService) EnQueue(sendSmsDto *dto.SendSms) error {
	message := url.QueryEscape(sendSmsDto.Message)
	sendSmsDto.Message = message
	queue, err := NewQueue()

	if err != nil {
		log.Error(err.Error())
		return err
	}

	msgAttrs := map[string]*sqs.MessageAttributeValue{
		"Number": {
			DataType:    aws.String("String"),
			StringValue: aws.String(sendSmsDto.Number),
		},
	}
	err = queue.SendMessage(config.Params.SmsQueueName, "sendSms", message, msgAttrs)

	if err != nil {
		return err
	}

	return nil
}

func (svc *SmsService) Send(msg *sqs.Message) {
	externalSvcUrl := config.Params.SmsExternalServiceUrl
	msgBody := *msg.Body
	number := *msg.MessageAttributes["Number"].StringValue
	urlWithMsg := fmt.Sprintf("%s%s&Message=%s", externalSvcUrl, number, msgBody)
	res, err := http.Get(urlWithMsg)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to sent: %s", err.Error())
		log.Error(errMsg)
		return
	}

	var xmlResp struct {
		XMLName      xml.Name `xml:"ArrayOfServiceClass"`
		Text         string   `xml:",chardata"`
		ServiceClass struct {
			Text          string `xml:",chardata"`
			MessageId     string `xml:"MessageId"`
			Status        string `xml:"Status"`
			StatusText    string `xml:"StatusText"`
			ErrorCode     string `xml:"ErrorCode"`
			ErrorText     string `xml:"ErrorText"`
			SMSCount      string `xml:"SMSCount"`
			CurrentCredit string `xml:"CurrentCredit"`
		} `xml:"ServiceClass"`
	}

	data, err := io.ReadAll(res.Body)
	_ = res.Body.Close()

	if err != nil {
		log.Error(err)
		return
	}

	_ = xml.NewDecoder(bytes.NewReader(data)).Decode(&xmlResp)

	if xmlResp.ServiceClass.StatusText == "success" {
		log.WithFields(log.Fields{
			"number":           number,
			"urlWithMsg":       urlWithMsg,
			"gateway_response": string(data),
		}).Info("Delivery success.")

		q, _ := NewQueue()
		q.DeleteMessage(msg)
	} else {
		log.WithFields(log.Fields{
			"number":           number,
			"urlWithMsg":       urlWithMsg,
			"gateway_response": string(data),
		}).Error("Delivery failed.")
	}
}

func (svc *SmsService) ReceiveMessagePeriodic(queueName string) {
	q, err := NewQueue()

	if err != nil {
		log.Panic(err.Error())
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		messages := q.ReceiveMessage(queueName)

		for _, msg := range messages {
			go svc.Send(msg)
		}
	}
}
