package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mustafaa-dev/Go-Distance-Tracker-Microservices/types"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"time"
)

func main() {
	e := echo.New()
	e.GET("/:id", GetOTU)
	e.Logger.Fatal(e.Start(":8000"))
}

func GetOTU(c echo.Context) error {
	id := c.Param("id")
	d, _ := GetFromNats("otu-" + id)
	log.Println(d)
	return c.String(http.StatusOK, "Hello, OTU,"+id)
}

func GetFromNats(topic string) (*types.OTU, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	sub, err := nc.SubscribeSync(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	// Use a timeout context to prevent blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var data *types.OTU
	err = json.Unmarshal(msg.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	log.Println(msg)

	return data, nil
}
