package main

import (
    "calls"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

const (
	port    = ":3000"
	appName = "DEV NerfTurret Cambreur College A005"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	app.Use("/ws", calls.WsUpgrade)

	app.Get("/ws/:id", websocket.New(calls.WsInit))

	app.Get("/send/:data", calls.WsSendData)

	log.Fatal(app.Listen(port))
}
