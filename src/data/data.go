package data

import (
	"encoding/json"
	"log"
    "os"
    "fmt"
)

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type TurretPos struct {
    X float64
    Y float64
    Z float64
}

func (c *Coordinate) ToString() string {
    return fmt.Sprintf("%f;%f", c.X, c.Y)
}

var filepath string

func SetDataLocation(location string) {
    filepath = location
}

func FetchPcData(pcId int) (c Coordinate, e error) {
    bodyIn, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
        e = err
        return
	}

	var coordinates []Coordinate
	err = json.Unmarshal(bodyIn, &coordinates)
	if err != nil {
		log.Fatal(err)
	}
    c = coordinates[pcId-1]
    e = nil

    return
}

