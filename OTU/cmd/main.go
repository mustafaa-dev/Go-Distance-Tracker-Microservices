package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mustafaa-dev/Go-Distance-Tracker-Microservices/types"
	"log"
	"math/rand"
	"time"
)

const (
	Interval = 1 * time.Second
	WsURL    = "ws://localhost:3000/ws"
)

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
	defer conn.Close()
	for {
		otu.Coords = types.Coords{
			Lat: GenerateCord(),
			Lon: GenerateCord(),
		}
		if err := SendData(otu, conn); err != nil {
			log.Println("Error sending data:", err)
			break
		}
		time.Sleep(Interval)
	}

}

func GenerateCord() float64 {
	x := float64(rand.Intn(180))
	return x + rand.Float64()
}

func SendData(otu *types.OTU, conn *websocket.Conn) error {
	log.Println("Sending data to websocket:", otu)
	err := conn.WriteJSON(otu)
	if err != nil {
		log.Println("Error writing JSON to websocket:", err)
		return err
	}
	return nil
}
