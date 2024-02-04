package main

import (
    "calls"
    "github.com/NerfTurret/ini-parser"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
    "os"
    "fmt"
    "errors"
)

const configPathDefault = "../config.ini"

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

	app := fiber.New(fiber.Config{
		AppName: config["app.name"],
	})

	app.Use("/ws", calls.WsUpgrade)

	app.Get("/ws/:id", websocket.New(calls.WsInit))

	app.Get("/send/:data", calls.WsSendData)

    app.Get("/select/:id", calls.SelectComputerById)

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
