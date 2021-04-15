package service

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

	err = queue.SendMessage(config.Params.SmsQueueName, *sendSmsDto)

	if err != nil {
		return err
	}

	return nil
}

func (svc *SmsService) Send(msg *sqs.Message) {
	externalSrvcUrl := config.Params.SmsExternalServiceUrl
	msgBody := *msg.Body
	nmbr := *msg.MessageAttributes["Number"].StringValue
	url := fmt.Sprintf("%s%s&Message=%s", externalSrvcUrl, nmbr, msgBody)
	res, err := http.Get(url)

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
	res.Body.Close()

	if err != nil {
		log.Error(err)
		return
	}

	xml.NewDecoder(bytes.NewReader(data)).Decode(&xmlResp)

	if xmlResp.ServiceClass.StatusText == "success" {
		log.WithFields(log.Fields{
			"number":           nmbr,
			"url":              url,
			"gateway_response": string(data),
		}).Info("Delivery success.")

		q, _ := NewQueue()
		q.DeleteMessage(msg)
	} else {
		log.WithFields(log.Fields{
			"number":           nmbr,
			"url":              url,
			"gateway_response": string(data),
		}).Error("Delivery failed.")
	}
}
