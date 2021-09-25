package handler

import (
	"errors"
	"github.com/emamulandalib/airbringr-notification/dto"
	"github.com/emamulandalib/airbringr-notification/response"
	"github.com/emamulandalib/airbringr-notification/service"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SendEmail(c *fiber.Ctx) error {
	sendEmailDto := new(dto.SendEmail)

	if err := c.BodyParser(sendEmailDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.BodyParseFailedErrorMsg,
			Errors:  errors.New(response.BodyParseFailedErrorMsg),
		})
	}

	err := sendEmailDto.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.ValidationFailedMesaage,
			Errors:  err,
		})
	}

	svc := service.EmailService{}
	err = svc.EnQueue(sendEmailDto)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: "System failed to send the email.",
			Errors:  errors.New("system failed to send the email"),
		})
	}

	return c.JSON(response.Payload{
		Message: "Email will be delivered soon",
	})
}
