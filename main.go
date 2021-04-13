package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emamulandalib/airbringr-notification/handler"
	"github.com/emamulandalib/airbringr-notification/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const idleTimeout = 5 * time.Second

func main() {
	// setup app
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	// setup middlewares
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Dhaka",
	}))

	//routes
	app.Get("/", handler.Home)
	route.V1(app)

	// 404
	app.Use(handler.NotFound)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen("0.0.0.0:8080"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here

	fmt.Println("Fiber was successful shutdown.")
}
