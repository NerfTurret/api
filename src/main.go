package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

const (
	port    = ":3000"
	appName = "DEV NerfTurret Cambreur College A005"
)

var connections = make(map[*websocket.Conn]bool)

func main() {
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	//app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
	//	//* DEBUG
	//	log.Println(c.Locals("allowed"))
	//	log.Println(c.Params("id"))
	//	log.Println(c.Query("v"))
	//	log.Println(c.Cookies("session"))
	//
	//	if _, e := connections[c]; e {
	//		go func() {
	//			log.Println("Attempted duplicate conn id")
	//			if err := c.WriteMessage(websocket.TextMessage, []byte("ID already in use")); err != nil {
	//				log.Println("write:", err)
	//			}
	//		}()
	//	} else {
	//		connections[c] = true
	//	}
	//
	//	go func() {
	//		if err := c.WriteMessage(websocket.TextMessage, []byte("Connection established with id: "+c.Params("id"))); err != nil {
	//			delete(connections, c)
	//		}
	//	}()
	//}))

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		connections[c] = true

		if err := c.WriteMessage(websocket.TextMessage, []byte("Connection established with id: "+c.Params("id"))); err != nil {
			delete(connections, c)
		}

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				delete(connections, c)
				break
			}
			log.Printf("recv: %s", msg)
		}
	}))

	app.Get("/send/:data", func(c *fiber.Ctx) error {
		data := c.Params("data")
		for conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
				log.Println("write:", err)
				delete(connections, conn)
			}
		}
		return c.SendString("Data sent to all clients")
	})

	log.Fatal(app.Listen(port))
}
