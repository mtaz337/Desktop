package route

import (
	"github.com/emamulandalib/airbringr-notification/handler"
	"github.com/gofiber/fiber/v2"
)

func V1(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Get("/", handler.Home)
}
