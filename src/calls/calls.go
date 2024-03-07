package calls

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
    "data"
    "fmt"
    "log"
    "math"
    "strconv"
)

var Connections = make(map[*websocket.Conn]bool)
var Turret data.TurretPos
var GlobalPcZ float64

func SetTurretPos(t data.TurretPos) {
    Turret = t
}

func SetGlobalPcZ(n float64) {
    GlobalPcZ = n
}

func WsUpgrade(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
        c.Locals("allowed", true)
        return c.Next()
    }
    return fiber.ErrUpgradeRequired
}


func WsInit(c *websocket.Conn) {
    Connections[c] = true

    if err := c.WriteMessage(websocket.TextMessage,
    []byte("Connection established with id: "+c.Params("id")));
    err != nil {
        delete(Connections, c)
    }

    for {
        _, msg, err := c.ReadMessage()
        if err != nil {
            log.Println("readerr:", err)
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
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.SendStatus(fiber.ErrBadRequest.Code)
    }
    pos, err := data.FetchPcData(id)
    if err != nil {
        return c.SendStatus(fiber.ErrBadRequest.Code)
    }
    theta := math.Atan((Turret.X - pos.X)/math.Abs(pos.Y - Turret.Y))
    phi := math.Atan((Turret.Z - GlobalPcZ)/math.Abs(pos.Y - Turret.Y))
    wsSendData(fmt.Sprintf("%f;%f", theta, phi))
    return c.SendStatus(200)
}

