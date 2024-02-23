package main

import (
    "github.com/gofiber/fiber/v2/middleware/cors"
    "calls"
    "data"
    "github.com/NerfTurret/ini-parser"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
    "os"
    "fmt"
    "errors"
    "strconv"
)

const configPathDefault = "../config.ini"
const noTurretPosErrorMsg = "Config file must define the turret's coordinates: x, y, and z"
const incorrectTurretPosErrorMsg = "Invalid turret position defined in config file"
var Turret data.TurretPos

func main() {
    configPath, err := handleCommandLineArguments()
    if err != nil {
        log.Fatal(err)
        return
    }
    if configPath == "" {
        configPath = configPathDefault
    }

    fmt.Println("Current path to config file: ", configPath)

    config := map[string]string{}
    ini.ParseFromFile(configPath, config)

    if v, ok := config["turret.x"]; ok == false {
        log.Fatal(noTurretPosErrorMsg)
    } else {
        if Turret.X, err = strconv.ParseFloat(v, 64); err != nil {
            log.Fatal(incorrectTurretPosErrorMsg)
        }
    }
    if v, ok := config["turret.y"]; ok == false {
        log.Fatal(noTurretPosErrorMsg)
    } else {
        if Turret.Y, err = strconv.ParseFloat(v, 64); err != nil {
            log.Fatal(incorrectTurretPosErrorMsg)
        }
    }
    if v, ok := config["turret.z"]; ok == false {
        log.Fatal(noTurretPosErrorMsg)
    } else {
        if Turret.Z, err = strconv.ParseFloat(v, 64); err != nil {
            log.Fatal(incorrectTurretPosErrorMsg)
        }
    }
    calls.SetTurretPos(Turret)

    if v, ok := config["global.pcZ"]; ok == false {
        log.Fatal("Config file must define a global PC z coordinate")
    } else {
        if c, err := strconv.ParseFloat(v, 64); err != nil {
            log.Fatal("Invalid global PC z coordinate")
        } else {
            calls.SetGlobalPcZ(c)
        }
    }

    data.SetDataLocation(config["config.data"])
    data.FetchPcData(1)

	app := fiber.New(fiber.Config{
    	AppName: config["app.name"],
	})
    app.Use(cors.New())

	app.Use("/ws", calls.WsUpgrade)
	app.Get("/ws/:id", websocket.New(calls.WsInit))

    if (config["config.openApi"] != "0") {
        app.Get("/send/:data", calls.WsSendData)
    }

    app.Get("/select/:id", calls.SelectComputerById)

    app.Static("/", "./pub")
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendFile("./pub/index.html")
    })

	log.Fatal(app.Listen(config["config.port"]))
}

// First ret val -> config.ini path
func handleCommandLineArguments() (string, error) {
    if !(len(os.Args) > 1) {
        return "", nil
    }
    if os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("argv 1 -> filepath config.ini; default: \"./config.ini\"")
        return "", errors.New("")
    }
    return os.Args[1], nil
}
