package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mustafaa-dev/Go-Distance-Tracker-Microservices/types"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

const (
	NatsURL = "nats://localhost:4222"
	Subject = "otu"
)

type OTUConn struct {
	Conn    *websocket.Conn
	OTUChan chan types.OTU
}

func NewOTUConn() *OTUConn {
	return &OTUConn{
		OTUChan: make(chan types.OTU, 10),
	}
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}
	log.Println("Connected to NATS & Waiting For Connections")
	defer nc.Close()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, nc)
	})
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request, nc *nats.Conn) {
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New Connection Established")

	newConn := NewOTUConn()
	newConn.Conn = conn

	go func() {
		for {
			OTU := &types.OTU{}
			err = newConn.Conn.ReadJSON(OTU)
			if err != nil {
				log.Println("Error reading JSON from websocket:", err)
				return
			}
			newConn.OTUChan <- *OTU
		}
	}()

	go func() {
		for otu := range newConn.OTUChan {
			log.Println("Received OTU:", otu)
			newConn.SendData(nc)
		}
	}()
}
func (h *OTUConn) SendData(nc *nats.Conn) error {
	for otu := range h.OTUChan {
		log.Println("Sending data to NATS:", otu)
		data, err := json.Marshal(otu)
		if err != nil {
			log.Println("Error marshaling OTU to JSON:", err)
			return err
		}
		err = nc.Publish(Subject, data)
		if err != nil {
			log.Println("Error publishing to NATS:", err)
			return err
		}
	}
	return nil
}
