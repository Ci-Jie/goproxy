package server

import (
	"fmt"
	"goproxy/controller"
	"goproxy/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

// Start ...
func Start() {
	app := fiber.New()
	app.Use(logger.New())

	controller.Init()
	storage.Init()

	app.Get("*", controller.Handler)
	app.Listen(fmt.Sprintf(":%s", viper.GetString("port")))
}
