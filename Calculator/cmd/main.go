package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mustafaa-dev/Go-Distance-Tracker-Microservices/types"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"sync"
)

const (
	Subject = "otu"
	Cost    = 0.5
)

type OTUStates struct {
	OTUFirstState *types.OTU
	OTULastState  *types.OTU
	Distance      float64
	Cash          float64
}

func NewOTUStates(otu *types.OTU) *OTUStates {
	return &OTUStates{
		OTUFirstState: otu,
		OTULastState:  otu,
	}
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}
	log.Println("Connected to NATS & Waiting For Messages")
	defer nc.Close()
	sub, _ := nc.SubscribeSync(Subject)

	otus := make(map[string]*OTUStates)
	var otusMutex sync.Mutex

	go func() {
		for {
			ctx := context.Background()
			msg, err := sub.NextMsgWithContext(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				continue
			}
			otu := types.OTU{}
			err = json.Unmarshal(msg.Data, &otu)
			otusMutex.Lock()
			if err != nil {
				log.Println("Error unmarshaling OTU:", err)
				continue
			}

			o, exists := otus[otu.OTUID]
			if !exists {
				log.Println("Receiving coords from:", otu.OTUID)

				o = NewOTUStates(&otu)
				otus[otu.OTUID] = o
			} else {
				o.OTULastState = &otu
			}
			otusMutex.Unlock()
			go serveOTU(o, nc)

		}
	}()
	select {}
}

func serveOTU(otuStates *OTUStates, nc *nats.Conn) {
	defer func() {
		if otuStates.OTULastState != nil {
			//log.Println("First received coords from:", otuStates.OTUFirstState.OTUID, "Coords:", otuStates.OTUFirstState.Coords.Lat, otuStates.OTUFirstState.Coords.Lon)
			//log.Println("Last received coords from:", otuStates.OTULastState.OTUID, "Coords:", otuStates.OTULastState.Coords.Lat, otuStates.OTULastState.Coords.Lon)
			d := distance(otuStates.OTUFirstState.Coords.Lat, otuStates.OTUFirstState.Coords.Lon, otuStates.OTULastState.Coords.Lat, otuStates.OTULastState.Coords.Lon, "K")
			log.Printf("Total Distance of %v :: %v :: %d EGP", otuStates.OTUFirstState.OTUID, d, int(d*Cost))
			otuStates.Cash = d * Cost
			otuStates.Distance = d
			go func() {
				err := SendData(otuStates, nc)
				if err != nil {

				}
			}()
		}
	}()
	//log.Println("Received coords from:", otu.OTUID, "Coords:", otu.Coords.Lat, otu.Coords.Lon)
	//log.Printf("Received coords from: %v - total distance : %v", otu.OTU.OTUID, distance(otu.OTU.Coords.Lat, otu.OTU.Coords.Lon, otu.OTU.Coords.Lat, otu.OTU.Coords.Lon))
}
func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func SendData(otu *OTUStates, nc *nats.Conn) error {
	log.Printf("Sending %v Coords %v KM :: %v EGP NATS Server", otu.OTUFirstState.OTUID, otu.Distance, otu.Cash)
	data, err := json.Marshal(otu)
	if err != nil {
		log.Println("Error marshaling OTU to JSON:", err)
		return err
	}
	err = nc.Publish(fmt.Sprintf("otu-%v", otu.OTUFirstState.OTUID), data)
	if err != nil {
		log.Println("Error publishing to NATS:", err)
		return err
	}
	return nil
}
