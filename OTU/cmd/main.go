package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mustafaa-dev/Go-Distance-Tracker-Microservices/types"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	Interval = 1 * time.Second
	WsURL    = "ws://localhost:3000/ws"
)

var lastLat, lastLon float64

func NewOTU() *types.OTU {
	return &types.OTU{
		OTUID: uuid.New().String(),
		Coords: types.Coords{
			Lat: 0,
			Lon: 0,
		},
	}
}

func main() {
	otu := NewOTU()
	conn, _, err := websocket.DefaultDialer.Dial(WsURL, nil)
	if err != nil {
		log.Println("Error dialing websocket:", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	for {
		lat, lon := GenerateCord()
		otu.Coords = types.Coords{
			Lat: lat,
			Lon: lon,
		}
		if err := SendData(otu, conn); err != nil {
			log.Println("Error sending data:", err)
			break
		}
		time.Sleep(Interval)
	}

}

func GenerateCord() (float64, float64) {
	// Generate a small random increment
	latIncrement := math.Abs(rand.Float64()-0.5) / 100
	lonIncrement := math.Abs(rand.Float64()-0.5) / 100

	// Add the increment to the last generated coordinate
	lastLat += latIncrement
	lastLon += lonIncrement

	// Ensure the coordinates are within the valid range
	if lastLat > 90 {
		lastLat = 90
	} else if lastLat < -90 {
		lastLat = -90
	}
	if lastLon > 180 {
		lastLon = 180
	} else if lastLon < -180 {
		lastLon = -180
	}

	return lastLat, lastLon
}
func SendData(otu *types.OTU, conn *websocket.Conn) error {
	log.Printf("Sending Coords :: <%v,%v> of %v to WS Server", otu.Coords.Lat, otu.Coords.Lon, otu.OTUID)
	err := conn.WriteJSON(otu)
	if err != nil {
		log.Println("Error writing JSON to websocket:", err)
		return err
	}
	return nil
}
