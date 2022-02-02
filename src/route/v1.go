package route

import (
	"github.com/emamulandalib/airbringr-notification/handler"
	"github.com/gofiber/fiber/v2"
)

func V1(server *fiber.App, handler *handler.Handler) {
	v1 := server.Group("/v1")
	v1.Get("/", handler.Home)
	v1.Post("/send-sms", handler.SendSms)
	v1.Post("/send-email", handler.SendEmail)
}
