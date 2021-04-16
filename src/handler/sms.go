package handler

import (
	"errors"

	"github.com/emamulandalib/airbringr-notification/dto"
	"github.com/emamulandalib/airbringr-notification/response"
	"github.com/emamulandalib/airbringr-notification/service"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SendSms(c *fiber.Ctx) error {
	sendSmsDto := new(dto.SendSms)

	if err := c.BodyParser(sendSmsDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.BodyParseFailedErrorMsg,
			Errors:  errors.New(response.BodyParseFailedErrorMsg),
		})
	}

	err := sendSmsDto.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.ValidationFailedMesaage,
			Errors:  err,
		})
	}

	svc := service.SmsService{}
	err = svc.EnQueue(sendSmsDto)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: "System failed to send the sms.",
			Errors:  errors.New("system failed to send the sms"),
		})
	}

	return c.JSON(response.Payload{
		Message: "Message will be sent soon.",
	})
}
