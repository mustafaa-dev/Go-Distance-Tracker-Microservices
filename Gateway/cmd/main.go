package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mustafaa-dev/Go-Distance-Tracker-Microservices/types"
	"github.com/nats-io/nats.go"
	"net/http"
	"time"
)

type OTUStates struct {
	OTUFirstState *types.OTU
	OTULastState  *types.OTU
	Distance      float64
	Cash          float64
}

func main() {
	e := echo.New()
	e.GET("/:id", GetOTU)
	e.Logger.Fatal(e.Start(":8000"))
}

func GetOTU(c echo.Context) error {
	id := c.Param("id")
	d, err := GetFromNats("otu-" + id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "OTU Disconnected"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"otuID": d.OTUFirstState.OTUID, "distance": fmt.Sprintf("%v KM", d.Distance), "cash": fmt.Sprintf("%v EGP", d.Cash), "initCoords": d.OTUFirstState.Coords, "lastCoords": d.OTULastState.Coords})
}

func GetFromNats(topic string) (*OTUStates, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	sub, err := nc.SubscribeSync(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var data *OTUStates
	err = json.Unmarshal(msg.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return data, nil
}
