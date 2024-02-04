package calls

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
    "log"
)

var Connections = make(map[*websocket.Conn]bool)

func WsUpgrade(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
        c.Locals("allowed", true)
        return c.Next()
    }
    return fiber.ErrUpgradeRequired
}


func WsInit(c *websocket.Conn) {
    log.Println(c.Locals("allowed"))
    log.Println(c.Params("id"))
    log.Println(c.Query("v"))
    log.Println(c.Cookies("session")) 

    Connections[c] = true

    if err := c.WriteMessage(websocket.TextMessage,
    []byte("Connection established with id: "+c.Params("id")));
    err != nil {
        delete(Connections, c)
    }

    for {
        _, msg, err := c.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            delete(Connections, c)
            break
        }
        log.Printf("recv: %s", msg)
    }
}

func WsSendData(c *fiber.Ctx) error {
    log.Println(c.Params("data"))
    if err := wsSendData(c.Params("data")); err != nil {
        log.Println(err)
        return nil
    }
    return c.SendStatus(200)
}

func wsSendData(data string) error {
    for conn := range Connections {
        if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
            log.Println("write:", err)
            delete(Connections, conn)
            return err
        }
    }
    return nil
}

func SelectComputerById(c *fiber.Ctx) error {
    log.Println("Received: " + c.Params("id"))
    wsSendData(c.Params("id"))
    return c.SendStatus(200)
}
