package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Создание экземпляра Fiber
	app := fiber.New()

	// Определение маршрута /hello
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	// Запуск сервера на порту 8080
	log.Fatal(app.Listen(":8080"))
}
